package relay

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dezh-tech/immortal/types/filter"
	"github.com/dezh-tech/immortal/types/message"
	"golang.org/x/net/websocket"
)

// TODO::: replace with https://github.com/coder/websocket.

type Relay struct {
	conns map[*websocket.Conn]map[string]filter.Filters
}

func NewRelay() *Relay {
	return &Relay{
		conns: make(map[*websocket.Conn]map[string]filter.Filters),
	}
}

func (r *Relay) Start() error {
	http.Handle("/ws", websocket.Handler(r.handleWS))
	err := http.ListenAndServe(":3000", nil) //nolint

	return err
}

func (r *Relay) handleWS(ws *websocket.Conn) {
	// TODO::: replace with logger.
	log.Printf("new connection: %s\n", ws.RemoteAddr())

	// TODO::: make it concurrent safe.
	r.conns[ws] = make(map[string]filter.Filters)

	r.readLoop(ws)
}

func (r *Relay) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			// TODO::: replace with logger.
			log.Printf("error in connection handling: %s\n", err)

			// TODO::: drop connection?
			continue
		}

		msg := message.ParseMessage(buf[:n])
		if msg == nil {
			_, _ = ws.Write(message.MakeNotice("error: can't parse message.")) // TODO::: should we check error?

			continue
		}

		// TODO::: replace with logger.
		log.Printf("received envelope: %s\n", msg.String())

		switch msg.Type() {
		case "REQ":
			go r.HandleReq(ws, msg)

		case "EVENT":
			go r.HandleEvent(ws, msg)

		case "CLOSE":
			go r.HandleClose(ws, msg)

		default:
			break
		}
	}
}

func (r *Relay) HandleReq(ws *websocket.Conn, m message.Message) {
	// TODO::: loadfrom database and sent in first query based on limit.
	// TODO::: return EOSE.
	// TODO::: use a concurrent safe map.

	msg, ok := m.(*message.Req)
	if !ok {
		_, _ = ws.Write(message.MakeNotice("error: can't parse REQ message")) // TODO::: should we check error?

		return
	}

	subs, ok := r.conns[ws]
	if !ok {
		_, _ = ws.Write(message.MakeNotice(fmt.Sprintf("error: can't find connection %s",
			ws.RemoteAddr()))) // TODO::: should we check error?

		return
	}

	subs[msg.SubscriptionID] = msg.Filters

	// TODO::: return EVENT messages.
}

func (r *Relay) HandleEvent(ws *websocket.Conn, m message.Message) {
	// TODO::: send events to be stored and proccessed.

	// can we ignore assertion check?
	msg, ok := m.(*message.Event)
	if !ok {
		okm := message.MakeOK(false,
			"",
			"error: can't parse the message.", // TODO::: make an error builder.
		)

		_, _ = ws.Write(okm) // TODO::: should we check error?

		return
	}

	if !msg.Event.IsValid() {
		okm := message.MakeOK(false,
			msg.SubscriptionID,
			"invalid: invalid id or sig.", // TODO::: make an error builder.
		)

		_, _ = ws.Write(okm) // TODO::: should we check error?

		return
	}

	_, _ = ws.Write(message.MakeOK(true, msg.SubscriptionID, "")) // TODO::: should we check error?

	// TODO::: any better way?
	for conn, subs := range r.conns {
		for id, filters := range subs {
			if !filters.Match(msg.Event) {
				continue
			}
			_, _ = conn.Write(message.MakeEvent(id, msg.Event)) // TODO::: should we check error?
		}
	}
}

func (r *Relay) HandleClose(ws *websocket.Conn, m message.Message) {
	msg, ok := m.(*message.Close)
	if !ok {
		_, _ = ws.Write(message.MakeNotice("error: can't parse CLOSE message")) // TODO::: should we check error?

		return
	}

	conn, ok := r.conns[ws]
	if !ok {
		_, _ = ws.Write(message.MakeNotice(fmt.Sprintf("error: can't find connection %s",
			ws.RemoteAddr()))) // TODO::: should we check error?

		return
	}

	delete(conn, msg.String())
	_, _ = ws.Write(message.MakeClosed(msg.String(), "ok: closed successfully")) // TODO::: should we check error?
}

// Stop shutdowns the relay gracefully.
func (r *Relay) Stop() error {
	for wsConn, subs := range r.conns {
		for id := range subs {
			_, _ = wsConn.Write(message.MakeClosed(id, "relay is stopping.")) // TODO::: should we check error?
		}
		_ = wsConn.Close() // TODO::: should we check error?
	}

	return nil
}

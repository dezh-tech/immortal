package relay

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/dezh-tech/immortal/types/filter"
	"github.com/dezh-tech/immortal/types/message"
	"golang.org/x/net/websocket"
)

// TODO::: replace with https://github.com/coder/websocket.
// TODO::: replace `log` with main logger.

// Relay represents a nostr relay which keeps track of client connections and handle them.
type Relay struct {
	conns     map[*websocket.Conn]map[string]filter.Filters
	connsLock sync.RWMutex
}

func NewRelay() *Relay {
	return &Relay{
		conns:     make(map[*websocket.Conn]map[string]filter.Filters),
		connsLock: sync.RWMutex{},
	}
}

// Start strats a new relay instance.
func (r *Relay) Start() error {
	http.Handle("/ws", websocket.Handler(r.handleWS))
	err := http.ListenAndServe(":3000", nil) //nolint

	return err
}

// handleWS is WebSocket handler.
func (r *Relay) handleWS(ws *websocket.Conn) {
	log.Printf("new connection: %s\n", ws.RemoteAddr())

	r.connsLock.Lock()
	r.conns[ws] = make(map[string]filter.Filters)
	r.connsLock.Unlock()

	r.readLoop(ws)
}

// readLoop reads incoming messages from a client and answer to them.
func (r *Relay) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Printf("error in connection handling: %s\n", err)

			continue
		}

		msg := message.ParseMessage(buf[:n])
		if msg == nil {
			_, _ = ws.Write(message.MakeNotice("error: can't parse message."))

			continue
		}

		log.Printf("received envelope: %s\n", msg.String())

		switch msg.Type() {
		case "REQ":
			go r.handleReq(ws, msg)

		case "EVENT":
			go r.handleEvent(ws, msg)

		case "CLOSE":
			r.handleClose(ws, msg)
		}
	}
}

// handleReq handles new incoming REQ messages from client.
func (r *Relay) handleReq(ws *websocket.Conn, m message.Message) {
	// TODO::: loadfrom database and sent in first query based on limit.
	// TODO::: return EOSE.
	// TODO::: return EVENT messages.

	msg, ok := m.(*message.Req)
	if !ok {
		_, _ = ws.Write(message.MakeNotice("error: can't parse REQ message."))

		return
	}

	r.connsLock.Lock()
	defer r.connsLock.Unlock()

	subs, ok := r.conns[ws]
	if !ok {
		_, _ = ws.Write(message.MakeNotice(fmt.Sprintf("error: can't find connection %s.",
			ws.RemoteAddr())))

		return
	}

	subs[msg.SubscriptionID] = msg.Filters
}

// handleEvent handles new incoming EVENT messages from client.
func (r *Relay) handleEvent(ws *websocket.Conn, m message.Message) {
	msg, ok := m.(*message.Event)
	if !ok {
		okm := message.MakeOK(false,
			"",
			"error: can't parse EVENT message.",
		)

		_, _ = ws.Write(okm)

		return
	}

	if !msg.Event.IsValid() {
		okm := message.MakeOK(false,
			msg.SubscriptionID,
			"invalid: id or sig is not correct.",
		)

		_, _ = ws.Write(okm)

		return
	}

	_, _ = ws.Write(message.MakeOK(true, msg.SubscriptionID, ""))

	for conn, subs := range r.conns {
		for id, filters := range subs {
			if !filters.Match(msg.Event) {
				return
			}
			_, _ = conn.Write(message.MakeEvent(id, msg.Event))
		}
	}
}

// handleClose handles new incoming CLOSE messages from client.
func (r *Relay) handleClose(ws *websocket.Conn, m message.Message) {
	msg, ok := m.(*message.Close)
	if !ok {
		_, _ = ws.Write(message.MakeNotice("error: can't parse CLOSE message."))

		return
	}

	r.connsLock.Lock()
	defer r.connsLock.Unlock()

	conn, ok := r.conns[ws]
	if !ok {
		_, _ = ws.Write(message.MakeNotice(fmt.Sprintf("error: can't find connection %s.",
			ws.RemoteAddr())))

		return
	}

	delete(conn, msg.String())
	_, _ = ws.Write(message.MakeClosed(msg.String(), "ok: closed successfully."))
}

// Stop shutdowns the relay gracefully.
func (r *Relay) Stop() error {
	r.connsLock.Lock()
	defer r.connsLock.Unlock()

	for wsConn, subs := range r.conns {
		// close all subscriptions.
		for id := range subs {
			_, _ = wsConn.Write(message.MakeClosed(id, "error: shutdowning the relay."))
		}

		// close connection.
		_ = wsConn.Close()
	}

	return nil
}

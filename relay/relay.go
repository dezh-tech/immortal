package relay

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/dezh-tech/immortal/types/envelope"
	"github.com/dezh-tech/immortal/types/filter"
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

		env := envelope.ParseMessage(buf[:n])
		if env == nil {
			continue
		}

		// TODO::: replace with logger.
		log.Printf("received envelope: %s\n", env.String())

		// TODO::: NIP-45, NIP-42.
		switch env.Label() {
		case "REQ":
			go r.HandleReq(ws, env)

		case "EVENT":
			go r.HandleEvent(ws, env)

		case "CLOSE":
			go r.HandleClose(ws, env) // should we pass env here?

		default:
			break
		}
	}
}

func (r *Relay) HandleReq(ws *websocket.Conn, e envelope.Envelope) {
	// TODO::: loadfrom database and sent in first query based on limit.
	// see: NIP-01.
	// TODO::: use a concurrent safe map.

	env, ok := e.(*envelope.ReqEnvelope)
	if !ok {
		return // TODO::: return EVENT message.
	}

	subs, ok := r.conns[ws]
	if !ok {
		return // TODO::: return EVENT message.
	}

	subs[env.SubscriptionID] = env.Filters

	// TODO::: return EVENT message.
}

func (r *Relay) HandleEvent(ws *websocket.Conn, e envelope.Envelope) {
	// TODO::: send events to be stored and proccessed.

	// can we ignore assertion check?
	env, ok := e.(*envelope.EventEnvelope)
	if !ok {
		res, _ := envelope.MakeOKEnvelope(false,
			"",
			"error: can't parse the message.", // TODO::: make an error builder.
		).MarshalJSON()

		_, _ = ws.Write(res)

		return
	}

	if !env.Event.IsValid() {
		res, _ := envelope.MakeOKEnvelope(false,
			env.SubscriptionID,
			"invalid: invalid _id_ or _sig_.", // TODO::: make an error builder.
		).MarshalJSON()

		_, _ = ws.Write(res)

		return
	}

	res, _ := envelope.MakeOKEnvelope(true, env.SubscriptionID, "").MarshalJSON()
	_, _ = ws.Write(res)

	// TODO::: any better way?
	for conn, subs := range r.conns {
		for id, filters := range subs {
			if !filters.Match(env.Event) {
				continue
			}
			resEnv := envelope.MakeEventEnvelope(id, env.Event)
			encodedResEnv, err := resEnv.MarshalJSON()
			if err != nil {
				continue
			}

			_, err = conn.Write(encodedResEnv)
			if err != nil {
				continue
			}
		}
	}
}

func (r *Relay) HandleClose(ws *websocket.Conn, e envelope.Envelope) {
	env, ok := e.(*envelope.CloseEnvelope)
	if !ok {
		//  TODO::: send NOTICE message.
		return
	}

	conn, ok := r.conns[ws]
	if !ok {
		// TODO::: send NOTICE message.
		return
	}

	delete(conn, env.String())
	// TODO::: what should we return here?
}

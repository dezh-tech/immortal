package server

import (
	"errors"
	"io"
	"log"
	"net/http"

	// TODO::: replace with https://github.com/coder/websocket.
	"github.com/dezh-tech/immortal/types/envelope"
	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) Start() error {
	http.Handle("/ws", websocket.Handler(s.handleWS))
	err := http.ListenAndServe(":3000", nil) //nolint

	return err
}

func (s *Server) handleWS(ws *websocket.Conn) {
	// TODO::: replace with logger.
	log.Printf("new connection: %s\n", ws.RemoteAddr())

	// TODO::: make it concurrent safe.
	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
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

		msg := buf[:n]
		env := envelope.ParseMessage(msg)

		// TODO::: replace with logger.
		log.Printf("received envelope: %s\n", env.String())
	}
}

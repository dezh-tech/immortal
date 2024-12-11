package server

import (
	"fmt"

	"github.com/dezh-tech/immortal/types/message"
	"github.com/dezh-tech/immortal/utils"
	"github.com/gorilla/websocket"
)

// handleReq handles new incoming REQ messages from client.
func (s *Server) handleReq(conn *websocket.Conn, m message.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer measureLatency(s.metrics.RequestLatency)()

	status := success
	defer func() {
		s.metrics.RequestsTotal.WithLabelValues(status).Inc()
	}()

	msg, ok := m.(*message.Req)
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice("error: can't parse REQ message."))

		status = parseFail

		return
	}

	client, ok := s.conns[conn]
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: can't find connection %s",
			conn.RemoteAddr())))

		status = serverFail

		return
	}

	if s.config.Limitation.AuthRequired && !*client.isKnown {
		client.challenge = utils.GenerateChallenge(10)
		authm := message.MakeAuth(client.challenge)

		closem := message.MakeClosed(
			msg.SubscriptionID,
			"auth-required: we can't serve notes to unauthenticated users",
		)

		_ = conn.WriteMessage(1, closem)

		_ = conn.WriteMessage(1, authm)
		status = authFail

		return
	}

	if len(msg.Filters) >= s.config.Limitation.MaxFilters {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: max limit of filters is: %d",
			s.config.Limitation.MaxFilters)))

		status = limitsFail

		return
	}

	if s.config.Limitation.MaxSubidLength <= len(msg.SubscriptionID) {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: max limit of sub id is: %d",
			s.config.Limitation.MaxSubidLength)))

		status = limitsFail

		return
	}

	if len(client.subs) >= s.config.Limitation.MaxSubscriptions {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: max limit of subs is: %d",
			s.config.Limitation.MaxSubscriptions)))

		status = limitsFail

		return
	}

	res, err := s.handlers.HandleReq(msg.Filters)
	if err != nil {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: can't process REQ message: %s", err.Error())))
		status = databaseFail

		return
	}

	for _, e := range res {
		msg := message.MakeEvent(msg.SubscriptionID, &e)
		_ = conn.WriteMessage(1, msg)
	}

	_ = conn.WriteMessage(1, message.MakeEOSE(msg.SubscriptionID))

	client.Lock()
	s.metrics.Subscriptions.Inc()
	client.subs[msg.SubscriptionID] = msg.Filters
	client.Unlock()
}

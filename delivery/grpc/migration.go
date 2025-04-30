package grpc

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	rpb "github.com/dezh-tech/immortal/delivery/grpc/gen"
	"github.com/dezh-tech/immortal/delivery/websocket"
	"github.com/dezh-tech/immortal/infrastructure/redis"
	"github.com/dezh-tech/immortal/pkg/logger"
	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"
)

type migrationServer struct {
	*Server
}

func newMigrationServer(server *Server) *migrationServer {
	return &migrationServer{
		Server: server,
	}
}

func (m *migrationServer) ExportEvents(rawFilter *rpb.Filter, stream rpb.Migration_ExportEventsServer) error {
	f, err := filter.Decode(rawFilter.Raw)
	if err != nil {
		return err
	}

	events, err := m.handler.HandleReq(f, rawFilter.Pubkey)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := stream.Send(EventToProto(&event)); err != nil {
			logger.Error("can't send event to manager", "error", err.Error(), "eventID", event.ID)

			continue
		}
	}

	return nil
}

// todo::: maybe we can enhance error handling?
// todo::: should we execute kind 5 and 62 here?
// todo::: how to deal with protected events?
func (m *migrationServer) ImportEvents(stream rpb.Migration_ImportEventsServer) error {
	for {
		e, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		evt, err := event.Decode(e.Raw)
		if err != nil {
			continue
		}

		if evt.Kind.IsEphemeral() ||
			evt.Kind == types.KindEventDeletionRequest ||
			evt.Kind == types.KindRightToVanish {
			continue
		}

		if !evt.IsValid(evt.GetRawID()) {
			continue
		}

		if !checkLimitations(m.redis, *m.keeper.WebsocketServer.GetLimitation(), evt) {
			continue
		}

		if err := m.handler.HandleEvent(evt); err != nil {
			continue
		}

		eID := evt.GetRawID()
		if err := m.redis.AddEventToBloom(eID[:]); err != nil {
			continue
		}
	}

	return stream.SendAndClose(&rpb.ImportEventResponse{
		Success: true,
		Message: "All events imported successfully",
	})
}

func EventToProto(e *event.Event) *rpb.Event {
	raw, _ := e.Encode()

	return &rpb.Event{
		Raw: raw,
	}
}

func checkLimitations(r *redis.Redis,
	limits websocket.Limitation, evt *event.Event,
) bool {
	eID := evt.GetRawID()
	if err := r.CheckAcceptability(limits.RestrictedWrites, eID[:], evt.PublicKey); err != nil {
		return false
	}

	expirationTag := evt.Tags.GetValue("expiration")

	if expirationTag != "" {
		expiration, err := strconv.ParseInt(expirationTag, 10, 64)
		if err != nil {
			return false
		}

		if time.Now().Unix() >= expiration {
			return false
		}

		if err := r.AddDelayedTask(websocket.ExpirationTaskListName,
			fmt.Sprintf("%s:%d", evt.ID, evt.Kind), time.Until(time.Unix(expiration, 0))); err != nil {
			return false
		}
	}

	if len(evt.Content) > int(limits.MaxContentLength) {
		return false
	}

	if evt.Difficulty() < int(limits.MinPowDifficulty) {
		return false
	}

	if len(evt.Tags) > int(limits.MaxEventTags) {
		return false
	}

	if evt.CreatedAt < limits.CreatedAtLowerLimit ||
		evt.CreatedAt > limits.CreatedAtUpperLimit {
		return false
	}

	return true
}

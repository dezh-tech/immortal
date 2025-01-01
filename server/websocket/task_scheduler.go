package websocket

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dezh-tech/immortal/types"
)

func (s *Server) checkExpiration() {
	for range time.Tick(time.Minute) {
		tasks, err := s.redis.GetReadyTasks("expiration_events")
		if err != nil {
			log.Println("error in checking expired events", err)
		}

		failedTasks := make([]string, 0)

		if len(tasks) != 0 {
			for _, job := range tasks {
				data := strings.Split(job, ":")

				if len(data) != 2 {
					continue
				}

				kind, err := strconv.Atoi(data[1])
				if err != nil {
					continue
				}

				if err := s.handler.DeleteByID(data[0],
					types.Kind(kind)); err != nil { //nolint
					failedTasks = append(failedTasks, job)
				}
			}
		}

		if len(failedTasks) != 0 {
			for _, fj := range failedTasks {
				if err := s.redis.AddDelayedTask("expiration_events",
					fj, time.Minute*10); err != nil {
					continue // todo::: retry then send log to manager.
				}
			}
		}
	}
}

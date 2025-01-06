package websocket

import (
	"strconv"
	"strings"
	"time"

	"github.com/dezh-tech/immortal/types"
)

const expirationTaskListName = "expiration_events"

func (s *Server) checkExpiration() {
	for range time.Tick(time.Minute) {
		tasks, err := s.redis.GetReadyTasks(expirationTaskListName)
		if err != nil {
			continue
		}

		failedTasks := make([]string, 0)

		if len(tasks) != 0 {
			for _, task := range tasks {
				data := strings.Split(task, ":")

				if len(data) != 2 {
					continue
				}

				kind, err := strconv.Atoi(data[1])
				if err != nil {
					continue
				}

				if err := s.handler.DeleteByID(data[0],
					types.Kind(kind)); err != nil { //nolint
					failedTasks = append(failedTasks, task)
				}
			}
		}

		if len(failedTasks) != 0 {
			for _, ft := range failedTasks {
				if err := s.redis.AddDelayedTask(expirationTaskListName,
					ft, time.Minute*10); err != nil {
					continue
				}
			}
		}
	}
}

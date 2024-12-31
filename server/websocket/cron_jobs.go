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
		jobs, err := s.redis.GetReadyJobs("expiration_events")
		if err != nil {
			log.Println("error in checking expired events", err)
		}

		failedJobs := make([]string, 0)

		if len(jobs) != 0 {
			log.Println("got jobs...", jobs)

			for _, job := range jobs {
				data := strings.Split(job, ":")

				if len(data) != 2 {
					continue
				}

				kind, err := strconv.Atoi(data[1])
				if err != nil {
					failedJobs = append(failedJobs, job)
				}

				if err := s.handler.DeleteByID(data[0],
					types.Kind(kind)); err != nil { //nolint
					failedJobs = append(failedJobs, job)
				}
			}
		}

		if len(failedJobs) != 0 {
			for _, fj := range failedJobs {
				if err := s.redis.AddDelayedJob("expiration_events",
					fj, time.Minute*10); err != nil {
					continue // todo::: retry then send log to manager.
				}
			}
		}
	}
}

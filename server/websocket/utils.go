package websocket

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	success      = "success"
	databaseFail = "db_fail"
	parseFail    = "parse_fail"
	limitsFail   = "limits_fail"
	serverFail   = "server_fail"
	invalidFail  = "invalid_fail"
)

func measureLatency(ht prometheus.Histogram) func() {
	start := time.Now()

	return func() {
		ht.Observe(time.Since(start).Seconds())
	}
}

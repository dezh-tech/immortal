package websocket

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func MeasureLatency(ht prometheus.Histogram) func() {
	start := time.Now()
	return func() {
		ht.Observe(time.Since(start).Seconds())
	}
}

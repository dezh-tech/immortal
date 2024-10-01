package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	EventCounter        prometheus.Counter
	RequestCounter      prometheus.Counter
	SubscriptionCounter prometheus.Gauge
	ConnectionCounter   prometheus.Gauge
	EventLaency         prometheus.Histogram
	RequestLatency      prometheus.Histogram
}

func New() *Metrics {
	eventC := promauto.NewCounter(prometheus.CounterOpts{
		Name: "event_counter",
		Help: "number of events sent to the relay.",
	})

	connC := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "connection_counter",
		Help: "number of open websocket connections.",
	})

	subC := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "subscription_counter",
		Help: "number of open subscription.",
	})

	reqC := promauto.NewCounter(prometheus.CounterOpts{
		Name: "request_counter",
		Help: "number of REQ messages sent to relay.",
	})

	eventL := promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "event_latency",
		Help: "time needed to request to an EVENT message.",
	})

	reqL := promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "requset_latency",
		Help: "time needed to request to a REQ message.",
	})

	return &Metrics{
		EventCounter:        eventC,
		ConnectionCounter:   connC,
		SubscriptionCounter: subC,
		RequestCounter:      reqC,
		EventLaency:         eventL,
		RequestLatency:      reqL,
	}
}

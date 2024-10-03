package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	EventsTotal    prometheus.Counter
	RequestsTotal  prometheus.Counter
	Subscriptions  prometheus.Gauge
	Connections    prometheus.Gauge
	EventLaency    prometheus.Histogram
	RequestLatency prometheus.Histogram
}

func New() *Metrics {
	eventsT := promauto.NewCounter(prometheus.CounterOpts{
		Name: "events_total",
		Help: "number of events sent to the relay.",
	})

	reqsT := promauto.NewCounter(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "number of REQ messages sent to relay.",
	})

	subs := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "subscriptions",
		Help: "number of open subscription.",
	})

	conns := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "connections",
		Help: "number of open websocket connections.",
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
		EventsTotal:    eventsT,
		Connections:    conns,
		Subscriptions:  subs,
		RequestsTotal:  reqsT,
		EventLaency:    eventL,
		RequestLatency: reqL,
	}
}

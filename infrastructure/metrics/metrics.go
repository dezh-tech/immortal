package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	EventsTotal    *prometheus.CounterVec
	RequestsTotal  *prometheus.CounterVec
	AuthsTotal     *prometheus.CounterVec
	MessagesTotal  prometheus.Counter
	Subscriptions  prometheus.Gauge
	Connections    prometheus.Gauge
	EventLatency   prometheus.Histogram
	RequestLatency prometheus.Histogram
	AuthLatency    prometheus.Histogram
}

func New() *Metrics {
	eventsT := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "events_total",
		Help: "number of events sent to the relay.",
	}, []string{"status"})

	reqsT := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "number of REQ messages sent to relay.",
	}, []string{"status"})

	authsT := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "auths_total",
		Help: "number of AUTH messages sent to relay.",
	}, []string{"status"})

	msgT := promauto.NewCounter(prometheus.CounterOpts{
		Name: "messages_total",
		Help: "number of messages received.",
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
		Name: "request_latency",
		Help: "time needed to request to a REQ message.",
	})

	authL := promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "auth_latency",
		Help: "time needed to request to a AUTH message.",
	})

	return &Metrics{
		EventsTotal:    eventsT,
		AuthsTotal:     authsT,
		Connections:    conns,
		MessagesTotal:  msgT,
		Subscriptions:  subs,
		RequestsTotal:  reqsT,
		EventLatency:   eventL,
		RequestLatency: reqL,
		AuthLatency:    authL,
	}
}

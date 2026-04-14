package metric

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "my_space"
	app_name  = "my_app"
)

type Metrics struct {
	requestCounter  prometheus.Counter
	responseCounter *prometheus.CounterVec
}

var metrics *Metrics

func Init(_ context.Context) error {
	metrics = &Metrics{
		requestCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      app_name + "_requests_total",
				Help:      "Количество запросов к серверу",
			},
		),
		responseCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      app_name + "_response_total",
				Help:      "Количество ответов от сервера",
			},
			[]string{"status", "method"},
		),
	}
	return nil
}

func IncRequestCounter() {
	metrics.requestCounter.Inc()
}
func IncResponseCounter(status, method string) {
	metrics.responseCounter.WithLabelValues(status, method).Inc()
}

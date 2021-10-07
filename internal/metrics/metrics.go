package metrics

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Metrics interface {
	IncHit(status int, method, path string)
	ObserveResponseTime(status int, method, path string, observeTime float64)
}

type prometheusMetrics struct {
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec
}

func NewPrometheusMetrics(address string, name string) (*prometheusMetrics, error) {
	var metric prometheusMetrics
	metric.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_hits_total",
	})

	if err := prometheus.Register(metric.HitsTotal); err != nil {
		return nil, errors.Wrap(err, "prometheus.Register")
	}

	metric.Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name + "_hits",
	}, []string{"status", "name", "path"})

	if err := prometheus.Register(metric.Hits); err != nil {
		return nil, errors.Wrap(err, "prometheus.Register")
	}

	metric.Times = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: name + "_times",
	}, []string{"status", "method", "path"})

	if err := prometheus.Register(metric.Times); err != nil {
		return nil, errors.Wrap(err, "prometheus.Register")
	}

	if err := prometheus.Register(collectors.NewBuildInfoCollector()); err != nil {
		return nil, errors.Wrap(err, "prometheus.Register")
	}
	return &metric, nil
}

func (m *prometheusMetrics) IncHit(status int, method, path string) {
	m.HitsTotal.Inc()
	m.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (m *prometheusMetrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	m.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}

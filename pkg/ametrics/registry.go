package ametrics

import (
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

type aMetricsRegistry struct {
	*prometheus.Registry
	constLabels []*io_prometheus_client.LabelPair
}

func newCustomMetricsRegistry(labels map[string]string) *aMetricsRegistry {
	c := &aMetricsRegistry{
		Registry: prometheus.NewRegistry(),
	}

	for k := range labels {
		name := k
		val := labels[k]
		c.constLabels = append(c.constLabels, &io_prometheus_client.LabelPair{
			Name:  &name,
			Value: &val,
		})
	}

	return c
}

func (g *aMetricsRegistry) Gather() ([]*io_prometheus_client.MetricFamily, error) {
	metricFamilies, err := g.Registry.Gather()

	for _, metricFamily := range metricFamilies {
		metrics := metricFamily.Metric
		for _, metric := range metrics {
			metric.Label = append(metric.Label, g.constLabels...)
		}
	}

	return metricFamilies, err
}

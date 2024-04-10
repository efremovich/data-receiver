package ametrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"net/http"
	"regexp"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var regex = regexp.MustCompile("[0-9]+")

type Middleware interface {
	WrapHandlerNetHttp(handlerName string, handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc
	MetricHandler() http.Handler
	AddNewCounterMetric(name string, desc string) (MetricCount, error)
	AddNewSummaryMetric(name string, desc string) (MetricBase, error)
	AddNewHistogramMetric(name string, desc string) (MetricBase, error)
	AddNewHistogramMetricWithCustomBucket(name string, desc string, buckets []float64) (MetricBase, error)
	AddNewGaugeMetric(name string, desc string) (MetricGauge, error)
	AddNewCounterMetricWithLabel(name string, desc string, labelsNames []string) (CounterVec, error)
	AddNewHistogramMetricWithLabel(name string, desc string, labelsNames []string) (HistogramVec, error)
	AddNewSummaryMetricWithLabel(name string, desc string, labelsNames []string) (SummaryVec, error)
	GetCounterMetric(name string) (MetricCount, bool)
	GetSummaryMetric(name string) (MetricBase, bool)
	GetHistogramMetric(name string) (MetricBase, bool)
	GetGaugeMetric(name string) (MetricGauge, bool)
	GetCounterMetricWithLabel(name string) (CounterVec, bool)
	GetHistogramMetricWithLabel(name string) (HistogramVec, bool)
	GetSummaryMetricWithLabel(name string) (SummaryVec, bool)
}

type MetricBase interface {
	prometheus.Observer
}

type MetricCount interface {
	Inc()
	Add(float64)
}

type MetricGauge interface {
	Set(float64)
	Inc()
	Dec()
	Add(float64)
	Sub(float64)
	SetToCurrentTime()
}

type hh struct {
	handler func(http.ResponseWriter, *http.Request)
}

func (m *middleware) newHH(handler func(http.ResponseWriter, *http.Request)) *hh {
	return &hh{
		handler: handler,
	}
}
func (h *hh) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler(w, r)
}

type metricList struct {
	metricsCount        map[string]prometheus.Counter
	metricsHistogram    map[string]prometheus.Histogram
	metricsGauge        map[string]prometheus.Gauge
	metricsSummary      map[string]prometheus.Summary
	metricsCountVec     map[string]CounterVec
	metricsHistogramVec map[string]HistogramVec
	metricsGaugeVec     map[string]prometheus.GaugeVec
	metricsSummaryVec   map[string]SummaryVec
	validateNames       map[string]struct{}
	sync.Mutex
}

type middleware struct {
	buckets  []float64
	registry aMetricsRegistry
	metricList
}

func (m *middleware) WrapHandlerNetHttp(handlerName string, handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	reg := prometheus.WrapRegistererWith(prometheus.Labels{"handler": handlerName}, &m.registry)

	requestsTotal := promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Количество запросов",
		}, []string{"method", "code"},
	)

	requestDuration := promauto.With(reg).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Время выполнения запроса",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 20, 30, 40, 50, 60},
		},
		[]string{"method", "code"},
	)

	requestSize := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_size_bytes",
			Help: "Размер запроса в байтах",
		},
		[]string{"method", "code"},
	)
	responseSize := promauto.With(reg).NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_response_size_bytes",
			Help: "Размер ответа в байтах",
		},
		[]string{"method", "code"},
	)

	base := promhttp.InstrumentHandlerCounter(
		requestsTotal,
		promhttp.InstrumentHandlerDuration(
			requestDuration,
			promhttp.InstrumentHandlerRequestSize(
				requestSize,
				promhttp.InstrumentHandlerResponseSize(
					responseSize,
					m.newHH(handler),
				),
			),
		),
	)

	return base.ServeHTTP
}
func (m *middleware) validateNameDuplicate(name string) bool {
	m.metricList.Lock()
	_, ok := m.validateNames[name]
	m.metricList.Unlock()

	if ok {
		return false
	}

	return true
}

func (m *middleware) validateName(name string) string {
	return regex.ReplaceAllString(name, "")
}

func (m *middleware) AddNewCounterMetric(name string, desc string) (MetricCount, error) {
	name = m.validateName(name)
	if !m.validateNameDuplicate(name) {
		return nil, fmt.Errorf("duplicate Metric %s", name)
	}

	m.metricList.Lock()
	counter := prometheus.NewCounter(prometheus.CounterOpts{Name: name, Help: desc})
	m.registry.MustRegister(counter)
	m.metricsCount[name] = counter
	m.validateNames[name] = struct{}{}
	m.metricList.Unlock()

	return counter, nil
}
func (m *middleware) AddNewCounterMetricWithLabel(name string, desc string, labelsNames []string) (CounterVec, error) {
	name = m.validateName(name)
	if !m.validateNameDuplicate(name) {
		return nil, fmt.Errorf("duplicate Metric %s", name)
	}
	counter, err := newCounterVec(name, desc, labelsNames, &m.registry)
	if err != nil {
		return nil, err
	}

	m.metricList.Lock()
	m.metricsCountVec[name] = counter
	m.validateNames[name] = struct{}{}
	m.metricList.Unlock()

	return counter, nil
}

func (m *middleware) GetCounterMetricWithLabel(name string) (CounterVec, bool) {
	name = m.validateName(name)
	m.metricList.Lock()
	metric, ok := m.metricsCountVec[name]
	m.metricList.Unlock()
	return metric, ok
}

func (m *middleware) GetCounterMetric(name string) (MetricCount, bool) {
	name = m.validateName(name)
	m.metricList.Lock()
	metric, ok := m.metricsCount[name]
	m.metricList.Unlock()
	return metric, ok
}

func (m *middleware) AddNewSummaryMetric(name string, desc string) (MetricBase, error) {
	name = m.validateName(name)
	if !m.validateNameDuplicate(name) {
		return nil, fmt.Errorf("duplicate Metric %s", name)
	}
	m.metricList.Lock()

	summary := prometheus.NewSummary(prometheus.SummaryOpts{Name: name, Help: desc})
	m.registry.MustRegister(summary)

	m.metricsSummary[name] = summary
	m.validateNames[name] = struct{}{}
	m.metricList.Unlock()

	return summary, nil
}

func (m *middleware) GetSummaryMetric(name string) (MetricBase, bool) {
	name = m.validateName(name)
	m.metricList.Lock()
	metric, ok := m.metricsSummary[name]
	m.metricList.Unlock()
	return metric, ok
}

func (m *middleware) AddNewHistogramMetric(name string, desc string) (MetricBase, error) {
	name = m.validateName(name)
	if !m.validateNameDuplicate(name) {
		return nil, fmt.Errorf("duplicate Metric %s", name)
	}
	m.metricList.Lock()

	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{Name: name, Help: desc, Buckets: m.buckets})
	m.registry.MustRegister(histogram)

	m.metricsHistogram[name] = histogram
	m.validateNames[name] = struct{}{}
	m.metricList.Unlock()

	return histogram, nil
}

func (m *middleware) AddNewHistogramMetricWithCustomBucket(name string, desc string, buckets []float64) (MetricBase, error) {
	name = m.validateName(name)
	if !m.validateNameDuplicate(name) {
		return nil, fmt.Errorf("duplicate Metric %s", name)
	}
	m.metricList.Lock()

	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{Name: name, Help: desc, Buckets: buckets})
	m.registry.MustRegister(histogram)

	m.metricsHistogram[name] = histogram
	m.validateNames[name] = struct{}{}
	m.metricList.Unlock()

	return histogram, nil
}

func (m *middleware) AddNewHistogramMetricWithLabel(name string, desc string, labelsNames []string) (HistogramVec, error) {
	name = m.validateName(name)
	if !m.validateNameDuplicate(name) {
		return nil, fmt.Errorf("duplicate Metric %s", name)
	}
	histogram, err := newHistogramVec(name, desc, labelsNames, &m.registry)
	if err != nil {
		return nil, err
	}

	m.metricList.Lock()
	m.metricsHistogramVec[name] = histogram
	m.validateNames[name] = struct{}{}
	m.metricList.Unlock()

	return histogram, nil
}

func (m *middleware) AddNewSummaryMetricWithLabel(name string, desc string, labelsNames []string) (SummaryVec, error) {
	name = m.validateName(name)
	if !m.validateNameDuplicate(name) {
		return nil, fmt.Errorf("duplicate Metric %s", name)
	}
	summary, err := newSummaryVec(name, desc, labelsNames, &m.registry)
	if err != nil {
		return nil, err
	}

	m.metricList.Lock()
	m.metricsSummaryVec[name] = summary
	m.validateNames[name] = struct{}{}
	m.metricList.Unlock()

	return summary, nil
}

func (m *middleware) GetHistogramMetricWithLabel(name string) (HistogramVec, bool) {
	name = m.validateName(name)
	m.metricList.Lock()
	metric, ok := m.metricsHistogramVec[name]
	m.metricList.Unlock()
	return metric, ok
}

func (m *middleware) GetSummaryMetricWithLabel(name string) (SummaryVec, bool) {
	name = m.validateName(name)
	m.metricList.Lock()
	metric, ok := m.metricsSummaryVec[name]
	m.metricList.Unlock()
	return metric, ok
}

func (m *middleware) GetHistogramMetric(name string) (MetricBase, bool) {
	name = m.validateName(name)
	m.metricList.Lock()
	metric, ok := m.metricsHistogram[name]
	m.metricList.Unlock()
	return metric, ok
}

func (m *middleware) AddNewGaugeMetric(name string, desc string) (MetricGauge, error) {
	name = m.validateName(name)
	if !m.validateNameDuplicate(name) {
		return nil, fmt.Errorf("duplicate Metric %s", name)
	}
	m.metricList.Lock()

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{Name: name, Help: desc})
	m.registry.MustRegister(gauge)

	m.metricsGauge[name] = gauge
	m.validateNames[name] = struct{}{}
	m.metricList.Unlock()

	return gauge, nil
}

func (m *middleware) GetGaugeMetric(name string) (MetricGauge, bool) {
	name = m.validateName(name)
	m.metricList.Lock()
	metric, ok := m.metricsGauge[name]
	m.metricList.Unlock()
	return metric, ok
}

func (m *middleware) MetricHandler() http.Handler {
	return promhttp.HandlerFor(&m.registry, promhttp.HandlerOpts{Registry: m.registry})
}

func New(serviceName, envelopment string) Middleware {
	registry := newCustomMetricsRegistry(prometheus.Labels{"serviceName": serviceName, "envelopment": envelopment})
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	buckets := []float64{0.1, 0.2, 0.35, 0.5, 0.75, 1, 1.25, 1.6, 2, 2.5, 3, 3.5, 4, 5, 6, 7.5}
	metrics := metricList{
		metricsCount:        map[string]prometheus.Counter{},
		metricsSummary:      map[string]prometheus.Summary{},
		metricsGauge:        map[string]prometheus.Gauge{},
		metricsHistogram:    map[string]prometheus.Histogram{},
		metricsCountVec:     map[string]CounterVec{},
		metricsHistogramVec: map[string]HistogramVec{},
		metricsSummaryVec:   map[string]SummaryVec{},
		validateNames:       map[string]struct{}{},
	}
	return &middleware{
		buckets:    buckets,
		registry:   *registry,
		metricList: metrics,
	}
}

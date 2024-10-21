package metrics

import (
	"net/http"
	"time"

	"github.com/efremovich/data-receiver/pkg/ametrics"
)

const (
	errorDescLabel  = "error_desc"
	apiCounterLabel = "method"
)

type Collector interface {
	ServeHTTP() http.Handler

	// Счетчик входящих заданий на создание служебных документов.
	IncServiceDocsTaskCounter()

	// Счетчик временных внутренних ошибок
	AddTemporaryError(errorDesc string)

	// Счетчик критических ошибок
	AddCriticalError(errorDesc string)

	// Время получение ответа от внешних источников.
	AddDocumentsAPIRequestTime(d time.Duration)
}

type metricsCollectorImplementation struct {
	metrics ametrics.Middleware

	serviceDocsTaskCounter           ametrics.MetricCount // Счетчик входящих заданий на создание служебных документов.
	apiMethodsMapCounter             ametrics.CounterVec  // Использование API сервиса.
	documentsAPIrequestTimeHistogram ametrics.MetricBase  // Время запросов в API документов.
	temporaryErrorMapCounter         ametrics.CounterVec  // Счетчик временных ошибок.
	criticalErrorsMapCounter         ametrics.CounterVec  // Счетчик критических ошибок.
}

func NewMetricCollector(serviceName string) (Collector, error) {
	metrics := ametrics.New("", serviceName)

	collector := metricsCollectorImplementation{
		metrics: metrics,
	}

	var err error

	collector.apiMethodsMapCounter, err = metrics.AddNewCounterMetricWithLabel("creator_api_methods_use_counter", "creator_api_methods_use_counter", []string{apiCounterLabel})
	if err != nil {
		return nil, err
	}

	collector.serviceDocsTaskCounter, err = metrics.AddNewCounterMetric("creator_task_counter", "creator_task_counter")
	if err != nil {
		return nil, err
	}

	collector.temporaryErrorMapCounter, err = metrics.AddNewCounterMetricWithLabel("creator_temporary_errors", "creator_temporary_errors", []string{errorDescLabel})
	if err != nil {
		return nil, err
	}

	collector.criticalErrorsMapCounter, err = metrics.AddNewCounterMetricWithLabel("creator_critical_errors", "creator_critical_errors", []string{errorDescLabel})
	if err != nil {
		return nil, err
	}

	bucket := []float64{0.02, 0.05, 0.1, 0.2, 0.35, 0.5, 0.75, 1, 1.25, 1.6, 2, 2.5, 3, 4}

	collector.documentsAPIrequestTimeHistogram, err = metrics.AddNewHistogramMetricWithCustomBucket("creator_documents_api_request_time", "creator_documents_api_request_time", bucket)
	if err != nil {
		return nil, err
	}

	return &collector, nil
}

func (m *metricsCollectorImplementation) AddDocumentsAPIRequestTime(d time.Duration) {
	if m.documentsAPIrequestTimeHistogram != nil {
		m.documentsAPIrequestTimeHistogram.Observe(d.Seconds())
	}
}

func (m *metricsCollectorImplementation) IncServiceDocsTaskCounter() {
	if m.serviceDocsTaskCounter != nil {
		m.serviceDocsTaskCounter.Inc()
	}
}

func (m *metricsCollectorImplementation) AddAPIMethodUse(method string) {
	if m.apiMethodsMapCounter != nil {
		_ = m.apiMethodsMapCounter.Inc(map[string]string{apiCounterLabel: method})
	}
}

func (m *metricsCollectorImplementation) AddTemporaryError(errorDesc string) {
	if m.temporaryErrorMapCounter != nil {
		_ = m.temporaryErrorMapCounter.Inc(map[string]string{errorDescLabel: errorDesc})
	}
}

func (m *metricsCollectorImplementation) AddCriticalError(errorDesc string) {
	if nil != m.criticalErrorsMapCounter {
		_ = m.criticalErrorsMapCounter.Inc(map[string]string{errorDescLabel: errorDesc})
	}
}

func (m *metricsCollectorImplementation) ServeHTTP() http.Handler {
	return m.metrics.MetricHandler()
}

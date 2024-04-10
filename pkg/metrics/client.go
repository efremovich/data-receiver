package metrics

import (
	"net/http"
	"time"

	"git.astralnalog.ru/utils/ametrics"
)

const (
	errorDescLabel  = "error_desc"
	apiCounterLabel = "method"
	operatorLabel   = "operator"
)

type Collector interface {
	ServeHTTP() http.Handler

	// Счетчик входящих ТП.
	IncTPCounter()

	// Использование методов API (кроме приёма ТП)
	AddAPIMethodUse(method string)

	// Количество ТП в данный момент в обработке + и -.
	IncCurrentTpInWork()
	DecCurrentTpInWork()

	// Время обработки ТП.
	AddTPProcessTime(d time.Duration)

	// Счетчик временных внутренних ошибок (500 ТРК)
	AddReceiveTPInternalError(errorDesc string)

	// Счетчик критических ошибок валидации ТП (400 ТРК)
	AddReceiveTPCriticalError(errorDesc string)

	// Время сохранения файла в сторадж.
	AddSaveStorageTime(d time.Duration)

	// Счётчик ТП по операторам.
	IncOperator(operator string)
}

type metricsCollectorImplementation struct {
	metrics ametrics.Middleware

	receiveTPCounter                ametrics.MetricCount // Просто счётчик приёма ТП.
	currentTPInWorksCounter         ametrics.MetricGauge // Количество ТП в данный момент в обработке.
	apiMethodsMapCounter            ametrics.CounterVec  // Использование API сервиса.
	receiveTPTimeHistogram          ametrics.MetricBase  // Время обработки ТП полностью.
	storageSaveTimeHistogram        ametrics.MetricBase  // Время сохранения ТП в сторадж.
	internalErrorMapCounter         ametrics.CounterVec  // Счетчик внутренних ошибок.
	validationErrorErrorsMapCounter ametrics.CounterVec  // Счетчик ошибок сломанных ТП.
	operatorsTPCounter              ametrics.CounterVec  // Счётчик ТП по оператоам.
}

func NewMetricCollector(serviceName string) (Collector, error) {
	metrics := ametrics.New("", serviceName)

	collector := metricsCollectorImplementation{
		metrics: metrics,
	}

	var err error

	collector.currentTPInWorksCounter, err = metrics.AddNewGaugeMetric("receiver_current_in_works", "receiver_current_in_works")
	if err != nil {
		return nil, err
	}

	collector.apiMethodsMapCounter, err = metrics.AddNewCounterMetricWithLabel("receiver_api_methods_use_counter", "receiver_api_methods_use_counter", []string{apiCounterLabel})
	if err != nil {
		return nil, err
	}

	bucket := []float64{0.1, 0.2, 0.35, 0.5, 0.75, 1, 1.25, 1.6, 2, 2.5, 3, 3.5, 4, 5, 6, 7.5}

	collector.receiveTPTimeHistogram, err = metrics.AddNewHistogramMetricWithCustomBucket("receiver_tp_process_time", "receiver_tp_process_time", bucket)
	if err != nil {
		return nil, err
	}

	collector.internalErrorMapCounter, err = metrics.AddNewCounterMetricWithLabel("receiver_temporary_errors", "receiver_temporary_errors", []string{errorDescLabel})
	if err != nil {
		return nil, err
	}

	collector.validationErrorErrorsMapCounter, err = metrics.AddNewCounterMetricWithLabel("receiver_critical_errors", "receiver_critical_errors", []string{errorDescLabel})
	if err != nil {
		return nil, err
	}

	bucket = []float64{0.02, 0.05, 0.1, 0.2, 0.35, 0.5, 0.75, 1, 1.25, 1.6, 2}

	collector.storageSaveTimeHistogram, err = metrics.AddNewHistogramMetricWithCustomBucket("receiver_storage_save_time", "receiver_storage_save_time", bucket)
	if err != nil {
		return nil, err
	}

	collector.operatorsTPCounter, err = metrics.AddNewCounterMetricWithLabel("receiver_operators", "receiver_operators", []string{operatorLabel})
	if err != nil {
		return nil, err
	}

	collector.receiveTPCounter, err = metrics.AddNewCounterMetric("receiver_receive_counter", "receiver_receive_counter")
	if err != nil {
		return nil, err
	}

	return &collector, nil
}

func (m *metricsCollectorImplementation) IncTPCounter() {
	m.receiveTPCounter.Inc()
}

func (m *metricsCollectorImplementation) AddAPIMethodUse(method string) {
	_ = m.apiMethodsMapCounter.Inc(map[string]string{apiCounterLabel: method})
}

func (m *metricsCollectorImplementation) DecCurrentTpInWork() {
	m.currentTPInWorksCounter.Inc()
}

func (m *metricsCollectorImplementation) IncCurrentTpInWork() {
	m.currentTPInWorksCounter.Dec()
}

func (m *metricsCollectorImplementation) AddTPProcessTime(d time.Duration) {
	m.receiveTPTimeHistogram.Observe(d.Seconds())
}

func (m *metricsCollectorImplementation) AddReceiveTPInternalError(errorDesc string) {
	_ = m.internalErrorMapCounter.Inc(map[string]string{errorDescLabel: errorDesc})
}

func (m *metricsCollectorImplementation) AddReceiveTPCriticalError(errorDesc string) {
	_ = m.validationErrorErrorsMapCounter.Inc(map[string]string{errorDescLabel: errorDesc})
}

func (m *metricsCollectorImplementation) AddSaveStorageTime(d time.Duration) {
	m.storageSaveTimeHistogram.Observe(d.Seconds())
}

func (m *metricsCollectorImplementation) IncOperator(operator string) {
	_ = m.operatorsTPCounter.Inc(map[string]string{operatorLabel: operator})
}

func (m *metricsCollectorImplementation) ServeHTTP() http.Handler {
	return m.metrics.MetricHandler()
}

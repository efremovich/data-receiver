package ametrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type HistogramVec interface {
	Observe(value float64, labelsWithValue map[string]string) error
}

type histogramVec struct {
	labelsCheck map[string]struct{}
	metric      *prometheus.HistogramVec
	sync.Mutex
}

func (h *histogramVec) checkLabels(labelsIn map[string]string) error {
	h.Lock()
	defer h.Unlock()

	if len(labelsIn) != len(h.labelsCheck) {
		return fmt.Errorf("не верное количество лейблов. Должно быть %d передано %d", len(labelsIn), len(h.labelsCheck))
	}
	for k, _ := range labelsIn {
		if _, ok := h.labelsCheck[k]; !ok {
			return fmt.Errorf("лейбл '%s' не входит в список инициализированных лейблов", k)
		}
	}

	return nil
}

func (h *histogramVec) Observe(value float64, labelsWithValue map[string]string) error {
	if err := h.checkLabels(labelsWithValue); err != nil {
		return err
	}

	h.metric.With(labelsWithValue).Observe(value)

	return nil

}

func newHistogramVec(name string, desc string, labels []string, reg *aMetricsRegistry) (HistogramVec, error) {
	labelsCheck := map[string]struct{}{}
	for _, v := range labels {
		if !isValidLabelName(v) {
			return nil, fmt.Errorf("лейбл должен начинаться с буквы и содержать символы [a-z A-Z 0-9 _ ]")
		}
		labelsCheck[v] = struct{}{}
	}

	hv := &histogramVec{
		metric:      prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: name, Help: desc}, labels),
		labelsCheck: labelsCheck,
	}

	reg.MustRegister(hv.metric)

	return hv, nil
}

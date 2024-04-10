package ametrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type CounterVec interface {
	Inc(labelsWithValue map[string]string) error
	Add(value float64, labelsWithValue map[string]string) error
}

type counterVec struct {
	labelsCheck map[string]struct{}
	metric      *prometheus.CounterVec
	sync.Mutex
}

func (c *counterVec) checkLabels(labelsIn map[string]string) error {
	c.Lock()
	defer c.Unlock()

	if len(labelsIn) != len(c.labelsCheck) {
		return fmt.Errorf("не верное количество лейблов. Должно быть %d передано %d", len(labelsIn), len(c.labelsCheck))
	}
	for k, _ := range labelsIn {
		if _, ok := c.labelsCheck[k]; !ok {
			return fmt.Errorf("лейбл '%s' не входит в список инициализированных лейблов", k)
		}
	}

	return nil
}

func (c *counterVec) Add(value float64, labelsWithValue map[string]string) error {
	if err := c.checkLabels(labelsWithValue); err != nil {
		return err
	}

	c.metric.With(labelsWithValue).Add(value)

	return nil

}
func (c *counterVec) Inc(labelsWithValue map[string]string) error {
	if err := c.checkLabels(labelsWithValue); err != nil {
		return err
	}

	c.metric.With(labelsWithValue).Inc()

	return nil

}

func newCounterVec(name string, desc string, labels []string, reg *aMetricsRegistry) (CounterVec, error) {
	labelsCheck := map[string]struct{}{}
	for _, v := range labels {
		if !isValidLabelName(v) {
			return nil, fmt.Errorf("лейбл должен начинаться с буквы и содержать символы [a-z A-Z 0-9 _ ]")
		}
		labelsCheck[v] = struct{}{}
	}

	cv := &counterVec{
		metric:      prometheus.NewCounterVec(prometheus.CounterOpts{Name: name, Help: desc}, labels),
		labelsCheck: labelsCheck,
	}

	reg.MustRegister(cv.metric)

	return cv, nil
}

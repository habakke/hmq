package metrics

import (
	"errors"
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
)

type Manager struct {
	prom *ginprom.Prometheus
}

const (
	MetricNumberOfMessages = "number_of_messages"
	MetricNumberOfClients  = "number_of_clients"
)

func (m *Manager) Init(e *gin.Engine) {
	if e == nil {
		log.Error("engine is nil, nothing to initialise")
		return
	}

	m.prom = ginprom.New(
		ginprom.Engine(e),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	m.prom.AddCustomGauge("test", "test metric", []string{})
	e.Use(m.prom.Instrument())
}

func (m *Manager) Add(name string, help string) error {
	if m.prom == nil {
		log.Error("metrics has not been initialised yet")
		return errors.New("metrics has not been initialised yet")
	}

	m.prom.AddCustomGauge(name, help, []string{})
	return nil
}

func (m *Manager) Inc(name string) error {
	if m.prom == nil {
		log.Error("metrics has not been initialised yet")
		return errors.New("metrics has not been initialised yet")
	}

	return m.prom.IncrementGaugeValue(name, []string{})
}

func (m *Manager) Dec(name string) error {
	if m.prom == nil {
		log.Error("metrics has not been initialised yet")
		return errors.New("metrics has not been initialised yet")
	}

	return m.prom.DecrementGaugeValue(name, []string{})
}

func (m *Manager) Set(name string, value float64) error {
	if m.prom == nil {
		log.Error("metrics has not been initialised yet")
		return errors.New("metrics has not been initialised yet")
	}

	return m.prom.SetGaugeValue(name, []string{}, value)
}

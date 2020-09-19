package monitor

import (
	"time"

	"github.com/mikemrm/go-venstar"
)

type Monitor struct {
	devices []*venstar.Device
}

func (m *Monitor) run(resultsChan chan *Results, errorsChan chan error) {
	var results *Results
	var err error
	for _, device := range m.devices {
		results, err = GetDeviceResults(device)
		if err != nil {
			errorsChan <- err
			continue
		}
		resultsChan <- results
	}
}

func (m *Monitor) IntervalMonitor(interval time.Duration, resultsChan chan *Results, errorsChan chan error, stopChan chan bool) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			m.run(resultsChan, errorsChan)
		}
	}
}

func (m *Monitor) Monitor(interval time.Duration, stopChan chan bool) (chan *Results, chan error) {
	resultsChan := make(chan *Results, 1)
	errorsChan := make(chan error, 1)
	go func() {
		defer close(resultsChan)
		defer close(errorsChan)
		m.IntervalMonitor(interval, resultsChan, errorsChan, stopChan)
	}()
	return resultsChan, errorsChan
}

type ResultsWriter interface {
	WriteResults(*Results) error
}

func New(devices ...*venstar.Device) *Monitor {
	return &Monitor{
		devices: devices,
	}
}

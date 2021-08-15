package prometheus

import (
	"fmt"
	"os"

	"go.mrm.dev/venstar-monitor"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	defaultVariableLabels = []string{"host", "type", "model", "name"}
)

type venstarCollector struct {
	monitor      *monitor.Monitor
	descriptions map[string]*prometheus.Desc
}

func newVenstarCollector(monitor *monitor.Monitor) *venstarCollector {
	collector := &venstarCollector{
		monitor: monitor,
	}
	collector.addDescriptions()
	return collector
}

func (c *venstarCollector) addDescriptions() {
	c.descriptions = make(map[string]*prometheus.Desc)
	c.addThermostatDescriptions()
}

func (c *venstarCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range c.descriptions {
		ch <- desc
	}
}

func (c *venstarCollector) Collect(ch chan<- prometheus.Metric) {
	resultsChan := make(chan *monitor.Results)
	errorsChan := make(chan error)
	doneChan := make(chan bool)
	go func() {
		defer func() {
			doneChan <- true
		}()
		for {
			select {
			case deviceResults, ok := <-resultsChan:
				if !ok {
					return
				}
				metrics, err := c.buildResultMetrics(deviceResults)
				if err != nil {
					fmt.Fprintln(os.Stderr, "[prometheus] Error building metric:", err)
					continue
				}
				for _, metric := range metrics {
					ch <- metric
				}
			case err, ok := <-errorsChan:
				if !ok {
					return
				}
				fmt.Fprintln(os.Stderr, "[prometheus] Error received:", err)
			}
		}
	}()
	c.monitor.Run(resultsChan, errorsChan)
	close(resultsChan)
	close(errorsChan)
	<-doneChan
}

func (c *venstarCollector) buildResultMetrics(results *monitor.Results) ([]prometheus.Metric, error) {
	tMetrics, err := c.buildThermostatMetrics(results)
	if err != nil {
		return nil, errors.Wrap(err, "building thermostat result metrics")
	}
	return tMetrics, nil
}

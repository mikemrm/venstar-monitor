package monitor

import (
	"github.com/mikemrm/go-venstar/thermostat"
)

type ThermostatResults struct {
	APIInfo   *thermostat.APIInfo
	QueryInfo *thermostat.QueryInfo
	Sensors   []*thermostat.Sensor
	Runtimes  []*thermostat.Runtime
	Alerts    []*thermostat.Alert
}

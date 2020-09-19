package monitor

import (
	"time"

	"github.com/mikemrm/go-venstar"
	"github.com/pkg/errors"
)

type Results struct {
	Timestamp  time.Time
	Device     *venstar.Device
	Thermostat *ThermostatResults
}

func GetDeviceResults(device *venstar.Device) (*Results, error) {
	results := &Results{
		Timestamp: time.Now(),
		Device:    device,
	}
	if device.Type == "thermostat" {
		apiinfo, err := device.Thermostat().GetAPIInfo()
		if err != nil {
			return nil, errors.Wrap(err, "loading api info")
		}
		queryinfo, err := device.Thermostat().GetQueryInfo()
		if err != nil {
			return nil, errors.Wrap(err, "loading query info")
		}
		sensors, err := device.Thermostat().GetQuerySensors()
		if err != nil {
			return nil, errors.Wrap(err, "loading query sensors")
		}
		runtimes, err := device.Thermostat().GetQueryRuntimes()
		if err != nil {
			return nil, errors.Wrap(err, "loading query runtimes")
		}
		alerts, err := device.Thermostat().GetQueryAlerts()
		if err != nil {
			return nil, errors.Wrap(err, "loading query alerts")
		}
		results.Thermostat = &ThermostatResults{
			APIInfo:   apiinfo,
			QueryInfo: queryinfo,
			Sensors:   sensors,
			Runtimes:  runtimes,
			Alerts:    alerts,
		}
	}
	return results, nil
}

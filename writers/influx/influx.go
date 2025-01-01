package influx

import (
	"fmt"
	"os"

	client "github.com/influxdata/influxdb1-client/v2"
	monitor "go.mrm.dev/venstar-monitor"
)

type Config struct {
	Addr string
	User string
	Pass string

	Database        string
	Measurement     string
	RetentionPolicy string
}

type InfluxWriter struct {
	client          client.Client
	Database        string
	Measurement     string
	RetentionPolicy string
}

func (w *InfluxWriter) WriteResults(results *monitor.Results) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:        w.Database,
		Precision:       "s",
		RetentionPolicy: w.RetentionPolicy,
	})
	if err != nil {
		return err
	}
	var point *client.Point
	if results.Device.Type == "thermostat" {
		point, err = client.NewPoint(w.Measurement, thermostatTags(results), thermostatFields(results), results.Timestamp)
		if err != nil {
			return err
		}
	} else {
		fmt.Fprintf(os.Stderr, "Unable to write influx points for unrecognized device type: %s\n", results.Device.Type)
	}
	bp.AddPoint(point)
	return w.client.Write(bp)
}

func thermostatTags(results *monitor.Results) map[string]string {
	tags := make(map[string]string)
	tags["device_type"] = results.Device.Type
	tags["device_host"] = results.Device.Address
	tResults := results.Thermostat
	if tResults.APIInfo != nil {
		if tResults.APIInfo.Model != "" {
			tags["system_model"] = fmt.Sprintf("%v", tResults.APIInfo.Model)
		}
		if tResults.APIInfo.Type != "" {
			tags["system_type"] = tResults.APIInfo.Type
		}
	}
	if tResults.QueryInfo != nil {
		if tResults.QueryInfo.Name != "" {
			tags["name"] = tResults.QueryInfo.Name
		}
		tags["state"] = tResults.QueryInfo.State.String()
	}
	return tags
}

func thermostatFields(results *monitor.Results) map[string]interface{} {
	fields := make(map[string]interface{})
	tResults := results.Thermostat
	if tResults.APIInfo != nil {
		fields["system_api_version"] = tResults.APIInfo.Version
		if tResults.APIInfo.Model != "" {
			fields["system_model"] = tResults.APIInfo.Model
		}
		if tResults.APIInfo.Firmware != "" {
			fields["system_firmware"] = tResults.APIInfo.Firmware
		}
		if tResults.APIInfo.Type != "" {
			fields["system_type"] = tResults.APIInfo.Type
		}
	}
	if tResults.QueryInfo != nil {
		fields["name"] = tResults.QueryInfo.Name
		fields["mode"] = int(tResults.QueryInfo.Mode)
		fields["state"] = int(tResults.QueryInfo.State)
		fields["fan"] = int(tResults.QueryInfo.Fan)
		fields["fan_state"] = int(tResults.QueryInfo.FanState)
		fields["temp_units"] = int(tResults.QueryInfo.TempUnits)
		fields["schedule"] = int(tResults.QueryInfo.Schedule)
		fields["schedule_part"] = int(tResults.QueryInfo.SchedulePart)
		fields["away"] = int(tResults.QueryInfo.Away)
		fields["holiday"] = int(tResults.QueryInfo.Holiday)
		fields["override"] = int(tResults.QueryInfo.Override)
		fields["override_remaining"] = int(tResults.QueryInfo.OverrideRemaining)
		fields["force_unoccupied"] = int(tResults.QueryInfo.ForceUnoccupied)
		fields["space_temp"] = tResults.QueryInfo.SpaceTemp
		fields["heat_temp"] = tResults.QueryInfo.HeatTemp
		fields["cool_temp"] = tResults.QueryInfo.CoolTemp
		fields["cool_temp_min"] = tResults.QueryInfo.CoolTempMin
		fields["cool_temp_max"] = tResults.QueryInfo.CoolTempMax
		fields["heat_temp_min"] = tResults.QueryInfo.HeatTempMin
		fields["heat_temp_max"] = tResults.QueryInfo.HeatTempMax
		fields["active_stage"] = tResults.QueryInfo.ActiveStage
		fields["humidity_enabled"] = int(tResults.QueryInfo.HumidityEnabled)
		fields["humidity"] = tResults.QueryInfo.Humidity
		fields["humidity_setpoint"] = tResults.QueryInfo.HumidifySetPoint
		fields["dehumidify_setpoint"] = tResults.QueryInfo.DehumidifySetPoint
		fields["setpoint_delta"] = tResults.QueryInfo.SetPointDelta
		fields["available_modes"] = int(tResults.QueryInfo.AvailableModes)
	}
	return fields
}

func NewWriter(config Config) (*InfluxWriter, error) {
	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Addr,
		Username: config.User,
		Password: config.Pass,
	})
	if err != nil {
		return nil, err
	}
	return &InfluxWriter{
		client:          client,
		Database:        config.Database,
		Measurement:     config.Measurement,
		RetentionPolicy: config.RetentionPolicy,
	}, nil
}

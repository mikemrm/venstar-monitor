package prometheus

import (
	"go.mrm.dev/venstar-monitor"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

func (c *venstarCollector) addThermostatDescriptions() {
	c.descriptions["ThermostatQueryInfoModeDesc"] = prometheus.NewDesc(
		"venstar_thermostat_mode",
		"Thermostat mode 0:off 1:heat 2:cool 3:auto",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoStateDesc"] = prometheus.NewDesc(
		"venstar_thermostat_state",
		"Thermostat current state 0:idle 1:heating 2:cooling 3:lockout 4:error",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoFanDesc"] = prometheus.NewDesc(
		"venstar_thermostat_fan",
		"Thermostat fan setting 0:auto 1:on",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoFanStateDesc"] = prometheus.NewDesc(
		"venstar_thermostat_fan_state",
		"Thermostat current fan state 0:off 1:on",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoTempUnitsDesc"] = prometheus.NewDesc(
		"venstar_thermostat_temp_units",
		"Thermostat temperature units 0:fahrenheit 1:celsius",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoScheduleDesc"] = prometheus.NewDesc(
		"venstar_thermostat_schedule",
		"Thermostat schedule 0:inactive 1:active",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoSchedulePartDesc"] = prometheus.NewDesc(
		"venstar_thermostat_schedule_part",
		"Thermostat schedule part 0:morning 1:day 2:evening 3:night 255:inactive",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoAwayDesc"] = prometheus.NewDesc(
		"venstar_thermostat_away",
		"Thermostat away 0:home 1:away",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoHolidayDesc"] = prometheus.NewDesc(
		"venstar_thermostat_holiday",
		"Thermostat holiday 0:not observing 1:observing",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoOverrideDesc"] = prometheus.NewDesc(
		"venstar_thermostat_override",
		"Thermostat override 0:off 1:on",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoOverrideRemainingDesc"] = prometheus.NewDesc(
		"venstar_thermostat_override_remaining_seconds",
		"Thermostat override remaining seconds",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoForceUnoccupiedDesc"] = prometheus.NewDesc(
		"venstar_thermostat_force_unoccupied",
		"Thermostat force unoccupied 0:off 1:on",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoSpaceTempDesc"] = prometheus.NewDesc(
		"venstar_thermostat_space_temp",
		"Thermostat current space temperature",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoHeatTempDesc"] = prometheus.NewDesc(
		"venstar_thermostat_heat_temp",
		"Thermostat set heat temperature",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoCoolTempDesc"] = prometheus.NewDesc(
		"venstar_thermostat_cool_temp",
		"Thermostat set cool temperature",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoCoolTempMinDesc"] = prometheus.NewDesc(
		"venstar_thermostat_cool_temp_min",
		"Thermostat set cool temperature min",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoCoolTempMaxDesc"] = prometheus.NewDesc(
		"venstar_thermostat_cool_temp_max",
		"Thermostat set cool temperature max",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoHeatTempMinDesc"] = prometheus.NewDesc(
		"venstar_thermostat_heat_temp_min",
		"Thermostat set heat temperature min",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoHeatTempMaxDesc"] = prometheus.NewDesc(
		"venstar_thermostat_heat_temp_max",
		"Thermostat set heat temperature max",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoActiveStageDesc"] = prometheus.NewDesc(
		"venstar_thermostat_active_stage",
		"Thermostat current heat/cool active stage",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoHumidityEnabledDesc"] = prometheus.NewDesc(
		"venstar_thermostat_humidity_enabled",
		"Thermostat humidity enabled 0:disabled 1:enabled",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoHumidityDesc"] = prometheus.NewDesc(
		"venstar_thermostat_humidity",
		"Thermostat humidity",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoHumidifySetPointDesc"] = prometheus.NewDesc(
		"venstar_thermostat_humidity_setpoint",
		"Thermostat humidity percentage setpoint",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoDehumidifySetPointDesc"] = prometheus.NewDesc(
		"venstar_thermostat_dehumidify_setpoint",
		"Thermostat dehumidify percentage setpoint",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoSetPointDeltaDesc"] = prometheus.NewDesc(
		"venstar_thermostat_setpoint_delta",
		"Thermostat setpoint delta",
		defaultVariableLabels,
		nil,
	)
	c.descriptions["ThermostatQueryInfoAvailableModesDesc"] = prometheus.NewDesc(
		"venstar_thermostat_available_modes",
		"Thermostat available modes 0:all 1:heat/cool 2:heat 3:cool",
		defaultVariableLabels,
		nil,
	)
}

func (c *venstarCollector) buildThermostatMetrics(results *monitor.Results) ([]prometheus.Metric, error) {
	if results == nil {
		return nil, nil
	}
	tResults := results.Thermostat
	labelValues := []string{
		results.Device.Address,
		tResults.APIInfo.Type,
		tResults.APIInfo.Model,
		tResults.QueryInfo.Name,
	}
	QueryInfoMode, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoModeDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.Mode),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoMode")
	}
	QueryInfoState, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoStateDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.State),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoState")
	}
	QueryInfoFan, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoFanDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.Fan),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoFan")
	}
	QueryInfoFanState, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoFanStateDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.FanState),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoFanState")
	}
	QueryInfoTempUnits, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoTempUnitsDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.TempUnits),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoTempUnits")
	}
	QueryInfoSchedule, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoScheduleDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.Schedule),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoSchedule")
	}
	QueryInfoSchedulePart, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoSchedulePartDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.SchedulePart),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoSchedulePart")
	}
	QueryInfoAway, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoAwayDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.Away),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoAway")
	}
	QueryInfoHoliday, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoHolidayDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.Holiday),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoHoliday")
	}
	QueryInfoOverride, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoOverrideDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.Override),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoOverride")
	}
	QueryInfoOverrideRemaining, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoOverrideRemainingDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.OverrideRemaining),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoOverrideRemaining")
	}
	QueryInfoForceUnoccupied, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoForceUnoccupiedDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.ForceUnoccupied),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoForceUnoccupied")
	}
	QueryInfoSpaceTemp, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoSpaceTempDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.SpaceTemp),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoSpaceTemp")
	}
	QueryInfoHeatTemp, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoHeatTempDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.HeatTemp),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoHeatTemp")
	}
	QueryInfoCoolTemp, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoCoolTempDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.CoolTemp),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoCoolTemp")
	}
	QueryInfoCoolTempMin, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoCoolTempMinDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.CoolTempMin),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoCoolTempMin")
	}
	QueryInfoCoolTempMax, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoCoolTempMaxDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.CoolTempMax),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoCoolTempMax")
	}
	QueryInfoHeatTempMin, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoHeatTempMinDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.HeatTempMin),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoHeatTempMin")
	}
	QueryInfoHeatTempMax, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoHeatTempMaxDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.HeatTempMax),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoHeatTempMax")
	}
	QueryInfoActiveStage, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoActiveStageDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.ActiveStage),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoActiveStage")
	}
	QueryInfoHumidityEnabled, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoHumidityEnabledDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.HumidityEnabled),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoHumidityEnabled")
	}
	QueryInfoHumidity, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoHumidityDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.Humidity),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoHumidity")
	}
	QueryInfoHumidifySetPoint, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoHumidifySetPointDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.HumidifySetPoint),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoHumidifySetPoint")
	}
	QueryInfoDehumidifySetPoint, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoDehumidifySetPointDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.DehumidifySetPoint),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoDehumidifySetPoint")
	}
	QueryInfoSetPointDelta, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoSetPointDeltaDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.SetPointDelta),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoSetPointDelta")
	}
	QueryInfoAvailableModes, err := prometheus.NewConstMetric(
		c.descriptions["ThermostatQueryInfoAvailableModesDesc"],
		prometheus.GaugeValue,
		float64(tResults.QueryInfo.AvailableModes),
		labelValues...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create metric for QueryInfoAvailableModes")
	}
	return []prometheus.Metric{
		QueryInfoMode,
		QueryInfoState,
		QueryInfoFan,
		QueryInfoFanState,
		QueryInfoTempUnits,
		QueryInfoSchedule,
		QueryInfoSchedulePart,
		QueryInfoAway,
		QueryInfoHoliday,
		QueryInfoOverride,
		QueryInfoOverrideRemaining,
		QueryInfoForceUnoccupied,
		QueryInfoSpaceTemp,
		QueryInfoHeatTemp,
		QueryInfoCoolTemp,
		QueryInfoCoolTempMin,
		QueryInfoCoolTempMax,
		QueryInfoHeatTempMin,
		QueryInfoHeatTempMax,
		QueryInfoActiveStage,
		QueryInfoHumidityEnabled,
		QueryInfoHumidity,
		QueryInfoHumidifySetPoint,
		QueryInfoDehumidifySetPoint,
		QueryInfoSetPointDelta,
		QueryInfoAvailableModes,
	}, nil
}

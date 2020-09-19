package main

import (
	"time"

	"github.com/mikemrm/venstar-monitor"
	"github.com/mikemrm/venstar-monitor/writers/influx"
	"github.com/mikemrm/venstar-monitor/writers/jsonPrinter"
)

var (
	outputFormats []string

	influxAddr            string
	influxUser            string
	influxPass            string
	influxDatabase        string
	influxMeasurement     string
	influxRetentionPolicy string
)

func init() {
	rootCmd.Flags().StringSliceVarP(&outputFormats, "output", "o", []string{"json"}, "One or more output formats")

	rootCmd.Flags().StringVar(&influxAddr, "influx.addr", "http://localhost:8086", "InfluxDB address")
	rootCmd.Flags().StringVar(&influxDatabase, "influx.database", "", "InfluxDB database (required)")
	rootCmd.Flags().StringVar(&influxMeasurement, "influx.measurement", "", "InfluxDB measurement (required)")
	rootCmd.Flags().StringVar(&influxRetentionPolicy, "influx.retention", "autogen", "InfluxDB retention policy")
	rootCmd.Flags().StringVar(&influxUser, "influx.user", "", "InfluxDB authentication username")
	rootCmd.Flags().StringVar(&influxPass, "influx.pass", "", "InfluxDB authentication password")
}

func loadOutputs() []monitor.ResultsWriter {
	var monitors []monitor.ResultsWriter
	for _, sWriter := range outputFormats {
		if sWriter == "json" {
			output := getJSONOutput()
			monitors = append(monitors, output)
		} else if sWriter == "influx" {
			output := getInfluxOutput()
			monitors = append(monitors, output)
		} else {
			fatal("Unknown output %s", sWriter)
		}
	}
	return monitors
}

func getJSONOutput() *jsonPrinter.JSONWriter {
	if monitoringInterval <= time.Duration(0) {
		fatal("json output requires an interval greater than 0")
	}
	writer, err := jsonPrinter.NewWriter()
	if err != nil {
		fatal("failed to create json writer: %s", err)
	}
	return writer
}

func getInfluxOutput() *influx.InfluxWriter {
	if influxAddr == "" {
		fatal("influx.addr must be defined")
	}
	if influxDatabase == "" {
		fatal("influx.db must be provided")
	}
	if influxMeasurement == "" {
		fatal("influx.measurement must be provided")
	}
	if influxRetentionPolicy == "" {
		fatal("influx.retention must be defined")
	}

	config := influx.Config{
		Addr:            influxAddr,
		Database:        influxDatabase,
		Measurement:     influxMeasurement,
		RetentionPolicy: influxRetentionPolicy,
		User:            influxUser,
		Pass:            influxPass,
	}

	writer, err := influx.NewWriter(config)
	if err != nil {
		fatal("failed to create influx writer: %s", err)
	}
	return writer
}

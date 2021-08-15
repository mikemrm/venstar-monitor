package main

import (
	"time"

	"go.mrm.dev/venstar-monitor"
	"go.mrm.dev/venstar-monitor/servers/prometheus"
	"go.mrm.dev/venstar-monitor/writers/influx"
	"go.mrm.dev/venstar-monitor/writers/jsonPrinter"
)

var (
	outputFormats []string

	influxAddr            string
	influxUser            string
	influxPass            string
	influxDatabase        string
	influxMeasurement     string
	influxRetentionPolicy string

	promListenAddr string
)

func init() {
	rootCmd.Flags().StringSliceVarP(&outputFormats, "output", "o", []string{"json"}, "One or more output formats")

	rootCmd.Flags().StringVar(&influxAddr, "influx.addr", "http://localhost:8086", "InfluxDB address")
	rootCmd.Flags().StringVar(&influxDatabase, "influx.database", "", "InfluxDB database (required)")
	rootCmd.Flags().StringVar(&influxMeasurement, "influx.measurement", "", "InfluxDB measurement (required)")
	rootCmd.Flags().StringVar(&influxRetentionPolicy, "influx.retention", "autogen", "InfluxDB retention policy")
	rootCmd.Flags().StringVar(&influxUser, "influx.user", "", "InfluxDB authentication username")
	rootCmd.Flags().StringVar(&influxPass, "influx.pass", "", "InfluxDB authentication password")

	rootCmd.Flags().StringVar(&promListenAddr, "prometheus.listen", ":9872", "Prometheus exporter listen address")
}

func loadOutputs(mon *monitor.Monitor) ([]monitor.ResultsWriter, map[string]monitor.Server) {
	var monitors []monitor.ResultsWriter
	servers := make(map[string]monitor.Server)
	for _, sWriter := range outputFormats {
		if sWriter == "json" {
			output := getJSONOutput()
			monitors = append(monitors, output)
		} else if sWriter == "influx" {
			output := getInfluxOutput()
			monitors = append(monitors, output)
		} else if sWriter == "prometheus" {
			servers["prometheus"] = getPrometheusOutput(mon)
		} else {
			fatal("Unknown output %s", sWriter)
		}
	}
	return monitors, servers
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

func getPrometheusOutput(mon *monitor.Monitor) monitor.Server {
	if promListenAddr == "" {
		fatal("prometheus.listen must be defined")
	}

	config := prometheus.Config{
		Monitor:    mon,
		ListenAddr: promListenAddr,
	}

	server, err := prometheus.NewServer(config)
	if err != nil {
		fatal("failed to create prometheus server: %s", err)
	}
	return server
}

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mikemrm/go-venstar"
	"github.com/mikemrm/venstar-monitor"
	"github.com/spf13/cobra"
)

var (
	thermostats        []string
	monitoringInterval time.Duration
)

func init() {
	rootCmd.Flags().StringSliceVarP(&thermostats, "thermostat", "t", nil, "One or more thermostat hosts to monitor")
	rootCmd.Flags().DurationVarP(&monitoringInterval, "interval", "i", 5*time.Second, "How frequent to poll devices for updates")
}

func fatal(s string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", values...)
	os.Exit(1)
}

func getDevices() []*venstar.Device {
	var devices []*venstar.Device
	for _, tstat := range thermostats {
		device, err := venstar.NewDevice("thermostat", tstat)
		if err != nil {
			fatal("Failed to create thermostat device (%s) %s", tstat, err)
		}
		devices = append(devices, device)
	}
	if len(devices) == 0 {
		fatal("No devices provided")
	}
	return devices
}

func runMonitor(cmd *cobra.Command, args []string) {
	monitors := loadOutputs()
	if len(monitors) != 0 {
		monitor := monitor.New(getDevices()...)
		stopChan := make(chan bool, 1)
		resultsChan, errorsChan := monitor.Monitor(monitoringInterval, stopChan)
		for {
			select {
			case results := <-resultsChan:
				for _, output := range monitors {
					err := output.WriteResults(results)
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error received writing results:", err)
					}
				}
			case err := <-errorsChan:
				fmt.Fprintln(os.Stderr, "Error received:", err)
			}
		}
	} else {
		fatal("No output formats specified")
	}
}

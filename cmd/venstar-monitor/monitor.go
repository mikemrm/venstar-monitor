package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.mrm.dev/venstar"
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
	wg := sync.WaitGroup{}
	mon := monitor.New(getDevices()...)
	monitors, servers := loadOutputs(mon)
	if len(monitors) == 0 && len(servers) == 0 {
		fatal("No output formats specified")
	}
	if len(monitors) != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stopChan := make(chan bool, 1)
			resultsChan, errorsChan := mon.Monitor(monitoringInterval, stopChan)
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
		}()
	}
	if len(servers) != 0 {
		for srvType, server := range servers {
			fmt.Fprintf(os.Stderr, "Starting %s server...\n", srvType)
			wg.Add(1)
			go func(t string, s monitor.Server) {
				defer wg.Done()
				err := s.Serve()
				if err != nil {
					fatal("%s server returned an error: %s", t, err)
				}
				fmt.Println(t, "DONE")
			}(srvType, server)
		}
		fmt.Fprintln(os.Stderr, "Servers started!")
	}
	wg.Wait()
}

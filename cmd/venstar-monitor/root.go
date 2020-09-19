package main

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "venstar-monitor",
		Short: "A service to monitor your Venstar devices",
		Run:   runMonitor,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

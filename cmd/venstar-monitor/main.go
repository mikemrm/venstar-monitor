package main

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/mikemrm/go-venstar"
// 	"github.com/mikemrm/venstar-monitor"
// 	"github.com/mikemrm/venstar-monitor/writers/jsonPrinter"
// 	// "github.com/mikemrm/venstar-monitor/monitoring/writers/influx"
// )

// func main() {
// }

// func main() {
// 	monitor := monitoring.New(os.Args[1])
// 	_ = monitor
// 	writer, err := influx.NewWriter(influx.Config{
// 		Addr:            "http://influx:8086",
// 		Database:        "homestats",
// 		Measurement:     "thermostat",
// 		RetentionPolicy: "autogen",
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	stopChan := make(chan bool, 1)
// 	resultsChan, errorsChan := monitor.Monitor(stopChan)
// 	for {
// 		select {
// 		case results := <-resultsChan:
// 			fmt.Println("Received results!")
// 			err := writer.WriteResults(results)
// 			if err != nil {
// 				fmt.Println("Error received writing results:", err)
// 			}
// 		case err := <-errorsChan:
// 			fmt.Println("Error received:", err)
// 		}
// 	}
// }

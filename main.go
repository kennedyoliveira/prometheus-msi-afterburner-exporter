package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kennedyoliveira/prometheus-msi-afterburner-exporter/monitor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	host            = flag.String("host", "127.0.0.1", "The host/ip of the machine running MSI Afterburner Remote Server")
	port            = flag.Int("port", 82, "The port of the machine running MSI Afterburner Remote Server")
	username        = flag.String("username", "MSIAfterburner", "Username of the MSI Afterburner Remote Server, by default it's MSIAfterburner, you should not touch it unless it actually changed")
	password        = flag.String("password", "17cc95b4017d496f82", "Password")
	help            = flag.Bool("help", false, "Prints the help information")
	updateInterval  = flag.Duration("update-interval", 2000*time.Millisecond, "Interval of polling the APIs for new monitoring data")
	listenAddress   = flag.String("listen-address", "0.0.0.0:8080", "The address and port that the prometheus metrics will be exposed")
	metricsEndpoint = flag.String("metrics-endpoint", "/metrics", "The path that the metrics endpoint will be available to be scrapped")
)

func main() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	log.Println("Starting...")
	host := fmt.Sprintf("%s:%d", *host, *port)

	log.Printf("Target host: %s", host)
	m := monitor.NewRemoteHardwareMonitor(*updateInterval, host, *username, *password)
	m.Start()

	http.Handle(*metricsEndpoint, promhttp.Handler())
	http.HandleFunc("/m/stop", func(writer http.ResponseWriter, request *http.Request) {
		m.Stop()
		writer.WriteHeader(204)
	})

	http.HandleFunc("/m/start", func(writer http.ResponseWriter, request *http.Request) {
		m.Start()
		writer.WriteHeader(204)
	})

	log.Printf("Listening at %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

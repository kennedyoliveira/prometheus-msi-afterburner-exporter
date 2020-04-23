package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/kennedyoliveira/prometheus-msi-afterburner-exporter/afterburner"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	host            = flag.String("host", "127.0.0.1", "The host/ip of the machine running MSI Afterburner Remote Server")
	port            = flag.Int("port", 82, "The port of the machine running MSI Afterburner Remote Server")
	username        = flag.String("username", "MSIAfterburner", "Username of the MSI Afterburner Remote Server, by default it's MSIAfterburner, you should not touch it unless it actually changed")
	password        = flag.String("password", "17cc95b4017d496f82", "Password")
	help            = flag.Bool("help", false, "Prints the help information")
	listenAddress   = flag.String("listen-address", "0.0.0.0:8080", "The address and port that the prometheus metrics will be exposed")
	metricsEndpoint = flag.String("metrics-endpoint", "/metrics", "The path that the metrics endpoint will be available to be scrapped")
)

func main() {
	version := GetVersionInfo()
	log.Printf("afterburner-exporter %s", version)
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	log.Println("Starting...")
	log.Printf("Target host: %s:%d", *host, *port)

	afterburnerClient := afterburner.NewAfterburnerClient(*host, *port, *username, *password)

	prometheus.MustRegister(NewVersionCollector(version))
	prometheus.MustRegister(NewAfterburnerCollector(afterburnerClient))

	http.Handle(*metricsEndpoint, promhttp.Handler())

	log.Printf("Listening at %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

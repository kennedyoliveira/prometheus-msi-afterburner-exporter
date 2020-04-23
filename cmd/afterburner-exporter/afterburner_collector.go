package main

import (
	"fmt"
	"log"

	"github.com/kennedyoliveira/prometheus-msi-afterburner-exporter/afterburner"
	"github.com/prometheus/client_golang/prometheus"
)

// collect and report metrics from afterburner
type AfterburnerCollector struct {
	client *afterburner.AfterburnerClient
}

func NewAfterburnerCollector(afterburnerClient *afterburner.AfterburnerClient) *AfterburnerCollector {
	return &AfterburnerCollector{client: afterburnerClient}
}

func (c *AfterburnerCollector) Describe(descs chan<- *prometheus.Desc) {
	// empty to allow for returning any metric, no check
}

func (c *AfterburnerCollector) Collect(metrics chan<- prometheus.Metric) {
	hardwareData, err := c.client.GetMonitorData()
	if err != nil {
		log.Printf("Could not get monitoring data: %v", err)
		return
	}

	for _, entry := range hardwareData.Data.Entries {
		metric, err := getMetric(&entry, &hardwareData.Gpus.Entries)
		if err != nil {
			switch err.(type) {
			case *BlackListedMetric:
				// ignore
				break
			default:
				log.Printf("Could not collect data from metric [%s]: %v", entry.LocalizedSourceName, err)
			}
		}

		if metric != nil {
			metrics <- *metric
		}
	}

	for _, gpu := range hardwareData.Gpus.Entries {
		metrics <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				"afterburner_gpu",
				"GPU information, the value is always 1 the information is on the tags",
				nil,
				prometheus.Labels{
					"name":   gpu.Device,
					"id":     gpu.GpuId,
					"bios":   gpu.BIOS,
					"driver": gpu.Driver,
					"family": gpu.Family,
				},
			),
			prometheus.GaugeValue,
			1,
		)
	}
}

func getMetric(metric *afterburner.HardwareMonitorEntry, gpus *[]afterburner.HardwareMonitorGpuEntry) (*prometheus.Metric, error) {
	name, labels, err := normalizeMetricName(metric, gpus)

	if err != nil {
		return nil, err
	}

	metricValue := normalizeMetricValue(metric.Data, metric.SourceUnits)

	promMetric := prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			fmt.Sprintf("afterburner_%s", name),
			//fmt.Sprintf("Metric extracted from %s with unit %s", metric.SourceName, metric.SourceUnits),
			"",
			nil,
			*labels,
		),
		prometheus.GaugeValue,
		metricValue,
	)

	return &promMetric, nil
}

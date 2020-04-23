package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Simple collector for showing the version and build information
// as metrics, the metric has a fixed name "afterburner_exporter_version" and the tags
// have the detailed information, like build number, tag, sha1 etc
type VersionCollector struct {
	versionDesc *prometheus.Desc
}

// creates and initialize a VersionCollector
func NewVersionCollector(info *VersionInfo) *VersionCollector {
	return &VersionCollector{
		versionDesc: prometheus.NewDesc(
			"afterburner_exporter_version",
			"Detailed information about the version and build",
			nil,
			prometheus.Labels{
				"branch":    info.Branch,
				"tag":       info.tag,
				"buildDate": info.buildDate,
				"sha1":      info.sha1,
			},
		)}
}

func (v *VersionCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- v.versionDesc
}

func (v *VersionCollector) Collect(metrics chan<- prometheus.Metric) {
	metrics <- prometheus.MustNewConstMetric(
		v.versionDesc,
		prometheus.GaugeValue,
		1,
	)
}

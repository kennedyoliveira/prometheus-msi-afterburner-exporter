package monitor

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"regexp"
	"strings"
)

type BlackListedMetric struct {
	MetricName string
}

func (b *BlackListedMetric) Error() string {
	return fmt.Sprintf("The metric [%s] is black listed", b.MetricName)
}

// regex for normalizing metrics
var (
	spaceRegex        = regexp.MustCompile("\\s+")
	invalidCharacters = regexp.MustCompile("[^a-zA-Z_ ][^a-zA-Z0-9_ ]*")
	cpuMetrics        = regexp.MustCompile("cpu.*?(\\d+)?.*")
)

// TODO load from config
var blackListRegex = []*regexp.Regexp{
	//regexp.MustCompile("cpu.*?usage.*"),
}

// TODO load from config
// regex to filter gpu metrics and assign them to
// the GPU
var gpuMetrics = []*regexp.Regexp{
	regexp.MustCompile(".*?gpu.*?"),
	regexp.MustCompile("(fb|vid|bus|memory) usage"),
	regexp.MustCompile("(core|memory) clock"),
	regexp.MustCompile("$power^"),
	regexp.MustCompile("fan (speed|tachometer)"),
	regexp.MustCompile("(temp|power|voltage|no load) limit"),
	// this are not gpu specific, it applies to the whole system
	//regexp.MustCompile("frame(rate|time).*"),
}

var metricUnits = map[string]string{
	"c":  "celsius",
	"ms": "millis",
	"w":  "watts",
	"%":  "percent",
	// this will be normalized to hertz
	"mhz": "hertz",
	"mb":  "bytes",
}

func normalizeMetricName(metric *HardwareMonitorEntry, gpus *[]HardwareMonitorGpuEntry) (string, *prometheus.Labels, error) {
	metricName := metric.LocalizedSourceName
	unit := metric.SourceUnits

	var name = strings.ToLower(strings.TrimSpace(metricName))

	for _, blRegex := range blackListRegex {
		if blRegex.MatchString(name) {
			return "", nil, &BlackListedMetric{
				MetricName: metricName,
			}
		}
	}

	var labels prometheus.Labels

	var suffix = ""

	lowerUnit := strings.ToLower(unit)
	if unitName, ok := metricUnits[lowerUnit]; ok {
		suffix = unitName
	} else {
		suffix = lowerUnit
	}

	labels = make(map[string]string)

	// check if it's a known metric
	if cpuMetrics.MatchString(name) {
		pieces := cpuMetrics.FindAllStringSubmatch(name, -1)

		if pieces[0][1] != "" {
			labels["core"] = pieces[0][1]
		} else {
			// reset the suffix to just total
			suffix = "total"
		}
	}

	for _, gpuRegex := range gpuMetrics {
		// if it's a know GPU metric
		if gpus != nil && gpuRegex.MatchString(name) {
			if len(*gpus) > metric.GpuIndex {
				gpu := (*gpus)[metric.GpuIndex]

				labels["gpu"] = gpu.Device
				labels["gpu_id"] = gpu.GpuId
				labels["gpu_bios"] = gpu.BIOS
				labels["gpu_driver"] = gpu.Driver
			} else {
				log.Printf("GPU index of %d but only %d gpus available", metric.GpuIndex+1, len(*gpus))
			}
		}
	}

	if suffix != "" {
		// we add the suffix in case there are any invalid character it will be removed
		name = name + "_" + suffix
	}

	// remove invalid characters
	name = invalidCharacters.ReplaceAllLiteralString(name, "")

	// transform spaces in _
	name = spaceRegex.ReplaceAllLiteralString(name, "_")

	return name, &labels, nil
}

func normalizeMetricValue(metricValue float64, metricUnit string) float64 {
	switch metricUnit {
	case "%":
		return metricValue / 100
	case "MB":
		return metricValue * 1000000
	case "MHz":
		return metricValue * 1000000
	default:
		return metricValue
	}
}

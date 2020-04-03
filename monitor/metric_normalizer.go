package monitor

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
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
var gpuMetrics = []*regexp.Regexp{}

var metricUnits = map[string]string{
	"c":  "celsius",
	"ms": "millis",
	"w":  "watts",
	"%":  "percent",
	// this will be normalized to hertz
	"mhz": "hertz",
}

func normalizeMetricName(metricName string, unit string) (string, *prometheus.Labels, error) {
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

	// check if it's a known metric
	if cpuMetrics.MatchString(name) {
		pieces := cpuMetrics.FindAllStringSubmatch(name, -1)
		labels = make(map[string]string)

		if pieces[0][1] != "" {
			labels["core"] = pieces[0][1]
		} else {
			// reset the suffix to just total
			suffix = "total"
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

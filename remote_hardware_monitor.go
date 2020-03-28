package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type RemoteHardwareMonitor struct {
	// how ofter the metrics will be scrapped
	updateInterval time.Duration

	// The host to be scrapped
	targetHost string

	ticker *time.Ticker

	// channel used to stop the monitoring
	done chan int
}

func NewRemoteHardwareMonitor(interval time.Duration, host string) *RemoteHardwareMonitor {
	return &RemoteHardwareMonitor{
		updateInterval: interval,
		targetHost:     host,
	}
}

// start the monitoring
func (hm *RemoteHardwareMonitor) Start() {
	log.Printf("Initializing the monitoring...")

	go func() {
		// clean up any previous ticker
		hm.Stop()

		log.Printf("Update interval: %v", hm.updateInterval)
		hm.ticker = time.NewTicker(hm.updateInterval)
		hm.done = make(chan int)

		for {
			select {
			case <-hm.done:
				log.Printf("Stopping the metrics polling...")
				return

			case <-hm.ticker.C:
				if err := pollMetrics(hm); err != nil {
					log.Printf("Failed to poll metrics: %v", err)
				}
			}
		}
	}()
}

// stop the monitoring
// it can be restarted
func (hm *RemoteHardwareMonitor) Stop() {
	if hm.ticker != nil {
		hm.ticker.Stop()
	}

	if hm.done != nil {
		close(hm.done)
	}
}

func pollMetrics(m *RemoteHardwareMonitor) error {
	url := fmt.Sprintf("http://%s/mahm", m.targetHost)

	httpClient := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return err
	}

	request.SetBasicAuth(*username, *password)

	requestsCounter.Inc()
	start := time.Now()
	response, err := httpClient.Do(request)

	requestsTime.Observe(time.Now().Sub(start).Seconds())

	if err != nil {
		failedRequestsCounter.Inc()
		log.Printf("Could not get information from %s: %v", url, err)
		return err
	}

	responsesCounter.WithLabelValues(string(response.StatusCode)).Inc()

	if response.StatusCode != 200 {
		return fmt.Errorf("Got invalid response from %v: %v", url, response.Status)
	}

	hardwareData, err := parseResponse(response.Body)

	if err != nil {
		return err
	}

	for _, entry := range hardwareData.Data.Entries {
		collector, err := getMetricCollector(entry.SourceName, entry.SourceUnits)
		if err != nil {
			switch err.(type) {
			case *BlackListedMetric:
				// ignore
				break
			default:
				log.Printf("Could not collect data from metric [%s]", entry.LocalizedSourceName)
			}
		}
		if collector != nil {
			// prometheus recommendation for percentage is to be between 0 and 1 not 0 and 100
			if entry.SourceUnits == "%" && entry.MaxLimit > 1 {
				(*collector).Set(entry.Data / 100.0)
			} else {
				(*collector).Set(entry.Data)
			}
		}
	}

	return nil
}

var (
	requestsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_collection_request_total",
		Help: "Amount of requests made to target host",
	})

	failedRequestsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_collection_request_fail_total",
		Help: "Amount of requests failed to the target host",
	})

	responsesCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_collection_response_total",
		Help: "Amount of responses from the target host by code",
	}, []string{"code"})

	requestsTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_collection_request_duration_seconds",
		Help:    "Histogram of the requests duration",
		Buckets: prometheus.DefBuckets,
	})
)

// regex for normalizing metrics
var (
	spaceRegex        = regexp.MustCompile("\\s+")
	invalidCharacters = regexp.MustCompile("[^a-zA-Z_ ][^a-zA-Z0-9_ ]*")
	cpuMetrics        = regexp.MustCompile("cpu.*?(\\d+)?.*")
)

var blackListRegex = []*regexp.Regexp{
	//regexp.MustCompile("cpu.*?usage.*"),
}

var metricUnits = map[string]string{
	"c":  "celsius",
	"ms": "millis",
	"w":  "watts",
	"%":  "percent",
}

type BlackListedMetric struct {
	MetricName string
}

func (b *BlackListedMetric) Error() string {
	return fmt.Sprintf("The metric [%s] is black listed", b.MetricName)
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
			suffix = "_total"
		}
	}

	// we add the suffix in case there are any invalid character it will be removed
	name = name + suffix

	// remove invalid characters
	name = invalidCharacters.ReplaceAllLiteralString(name, "")

	// transform spaces in _
	name = spaceRegex.ReplaceAllLiteralString(name, "_")

	return name, &labels, nil
}

func getMetricCollector(metricName string, metricUnit string) (*prometheus.Gauge, error) {
	name, labels, err := normalizeMetricName(metricName, metricUnit)

	if err != nil {
		return nil, err
	}

	labelNames := GetMapKeys(*labels)

	metricCollector := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: fmt.Sprintf("afterburner_%s", name),
	}, labelNames)

	if err := prometheus.Register(metricCollector); err != nil {
		if ex, ok := err.(prometheus.AlreadyRegisteredError); ok {
			existingGauge := ex.ExistingCollector.(*prometheus.GaugeVec)
			existingGaugeWithLabels := existingGauge.With(*labels)
			return &existingGaugeWithLabels, nil
		}

		return nil, err
	}

	gaugeWithLabels := metricCollector.With(*labels)
	return &gaugeWithLabels, nil
}

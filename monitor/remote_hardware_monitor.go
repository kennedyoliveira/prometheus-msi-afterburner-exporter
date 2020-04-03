package monitor

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"net/http"
	"time"
)

type RemoteHardwareMonitor struct {
	// how ofter the metrics will be scrapped
	updateInterval time.Duration

	// The host to be scrapped
	targetHost string

	// credentials to access the remote api
	username string
	password string

	// timer that will trigger the polling of data
	ticker *time.Ticker

	// channel used to stop the monitoring
	done chan int

	// gpus being monitored
	gpus []HardwareMonitorGpuEntry

	// amount of failure, for logging purposes
	failureCount int

	// http client for requests
	httpClient http.Client
}

func NewRemoteHardwareMonitor(interval time.Duration, host string, username string, password string) *RemoteHardwareMonitor {
	return &RemoteHardwareMonitor{
		updateInterval: interval,
		targetHost:     host,
		username:       username,
		password:       password,
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
		hm.httpClient = http.Client{
			Timeout: 10 * time.Second,
		}

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

	hm.gpus = nil
	hm.failureCount = 0
}

func queryMsiAfterburner(hm *RemoteHardwareMonitor) (*HardwareMonitor, error) {
	url := fmt.Sprintf("http://%s/mahm", hm.targetHost)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}

	request.SetBasicAuth(hm.username, hm.password)

	requestsCounter.Inc()
	start := time.Now()
	response, err := hm.httpClient.Do(request)

	requestsTime.Observe(time.Now().Sub(start).Seconds())

	if err != nil {
		failedRequestsCounter.Inc()
		hm.failureCount += 1
		log.Printf("Could not get information from %s: %v", url, err)
		return nil, err
	}

	defer response.Body.Close()

	responsesCounter.WithLabelValues(string(response.StatusCode)).Inc()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("got invalid response from %v: %v", url, response.Status)
	}

	hardwareData, err := parseResponse(response.Body)

	if err != nil {
		return nil, err
	}

	return hardwareData, nil
}

func collectMetrics(hardwareData *HardwareMonitor) {
	for _, entry := range hardwareData.Data.Entries {
		collector, err := getMetricCollector(&entry, &hardwareData.Gpus.Entries)
		if err != nil {
			switch err.(type) {
			case *BlackListedMetric:
				// ignore
				break
			default:
				log.Printf("Could not collect data from metric [%s]: %v", entry.LocalizedSourceName, err)
			}
		}
		if collector != nil {
			(*collector).Set(normalizeMetricValue(entry.Data, entry.SourceUnits))
		}
	}
}

func pollMetrics(hm *RemoteHardwareMonitor) error {
	hardwareData, err := queryMsiAfterburner(hm)

	if err != nil {
		return err
	}

	if (hm.gpus == nil || hm.failureCount >= 3) && len(hardwareData.Gpus.Entries) > 0 {
		hm.gpus = hardwareData.Gpus.Entries

		log.Printf("Gpus available in the host machine")

		for _, gpu := range hm.gpus {
			log.Printf("%s, Driver: %s, BIOS: %s, id: %s", gpu.Device, gpu.Driver, gpu.BIOS, gpu.GpuId)
		}
	}

	// reset the failure counts on success
	hm.failureCount = 0

	collectMetrics(hardwareData)

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

func getMetricCollector(metric *HardwareMonitorEntry, gpus *[]HardwareMonitorGpuEntry) (*prometheus.Gauge, error) {
	name, labels, err := normalizeMetricName(metric, gpus)

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

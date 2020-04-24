package afterburner

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type AfterburnerClient struct {
	host string
	port int

	username string
	password string

	httpClient *http.Client
}

func NewAfterburnerClient(host string, port int, username string, password string) *AfterburnerClient {
	return &AfterburnerClient{
		host:       host,
		port:       port,
		username:   username,
		password:   password,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// metrics for the scrap
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

func (c *AfterburnerClient) GetMonitorData() (*HardwareMonitor, error) {
	url := fmt.Sprintf("http://%s:%d/mahm", c.host, c.port)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}

	request.SetBasicAuth(c.username, c.password)

	requestsCounter.Inc()
	start := time.Now()
	response, err := c.httpClient.Do(request)

	requestsTime.Observe(time.Now().Sub(start).Seconds())

	if err != nil {
		failedRequestsCounter.Inc()
		log.Printf("Could not get information from %s: %v", url, err)
		return nil, err
	}

	defer response.Body.Close()

	responsesCounter.WithLabelValues(string(response.StatusCode)).Inc()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("got invalid response from %v: %v", url, response.Status)
	}

	hardwareData, err := ParseResponse(response.Body)

	if err != nil {
		return nil, err
	}

	return hardwareData, nil
}

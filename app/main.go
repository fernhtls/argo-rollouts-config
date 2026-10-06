package main

import (
	_ "embed"
	"fmt"
	"os"

	"encoding/json"
	"log"
	mrand "math/rand"
	"net/http"

	// "time"

	"gopkg.in/yaml.v3"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed config.yaml
var serverConfigBytes []byte

var serverConfig ServerCofing

type ServerCofing struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Version string `yaml:"version"`
}

// Metadata : metadata used for our request and response
type Metadata struct {
	Version string `json:"version"`
	Caller  string `json:"caller"`
}

// PredictRequest : struct for our request
type PredictRequest struct {
	Metadata Metadata  `json:"metadata"`
	Inputs   []float64 `json:"inputs"`
}

// PredictResponse : struct for the reponse
type PredictResponse struct {
	Metadata      Metadata  `json:"metadata"`
	Results       []float64 `json:"results"`
	ServerVersion string    `json:"server_version"`
}

// vars for our metrics / counter and duration
var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_requests_total",
			Help: "Total number of HTTP requests processed",
		},
		[]string{"method", "endpoint", "rollouts_pod_template_hash"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "rollouts_pod_template_hash"},
	)
)

// instrumentHandler : implements the metric instrumentation for a handler
func instrumentHandler(endpoint string, serverVersion string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(requestDuration.WithLabelValues(r.Method, endpoint, serverVersion))
		defer timer.ObserveDuration()
		requestCounter.WithLabelValues(r.Method, endpoint, os.Getenv("POD_TEMPLATE_HASH")).Inc()
		handler(w, r)
	}
}

// predict : handler function to our predict
func predict(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
	var pr PredictRequest
	if err := json.NewDecoder(req.Body).Decode(&pr); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()
	log.Printf("got request as : %v\n", pr)
	w.Header().Set("Content-type", "application/json")
	results := make([]float64, len(pr.Inputs))
	for i, v := range pr.Inputs {
		results[i] = v * float64(mrand.Intn(5))
	}
	// TODO: v2 will be adding the sleep here for requests to take longer
	// time.Sleep(time.Second*time.Duration(mrand.Intn(2)))
	prs := PredictResponse{
		Metadata:      pr.Metadata,
		Results:       results,
		ServerVersion: serverConfig.Version,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prs)
}

func init() {
	if err := yaml.Unmarshal(serverConfigBytes, &serverConfig); err != nil {
		log.Fatalf("error on parsing config: %v", err)
	}
}

func main() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		requestCounter,
		requestDuration,
	)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.Handle("/predict", instrumentHandler("predict", serverConfig.Version, predict))
	server := fmt.Sprintf("%s:%s", serverConfig.Server.Host, serverConfig.Server.Port)
	log.Printf("starting server %s with version %s ...\n", server, serverConfig.Version)
	if err := http.ListenAndServe(server, nil); err != nil {
		log.Fatalf("Server error: %s", err)
	}
}

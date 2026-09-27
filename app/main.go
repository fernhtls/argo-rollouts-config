package main

import (
	"encoding/json"
	"log"
	"net/http"
	// "time"
	mrand "math/rand"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metadata : metadata used for our request and response
type Metadata struct {
	Version string `json:"version"`
	Caller string `json:"caller"`
}

// PredictRequest : struct for our request
type PredictRequest struct {
	Metadata Metadata `json:"metadata"`
	Inputs []float64 `json:"inputs"`
}

// PredictResponse : struct for the reponse
type PredictResponse struct {
	Metadata Metadata `json:"metadata"`
	Results []float64 `json:"results"`
}

// vars for our metrics / counter and duration
var (
    requestCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "app_http_requests_total",
            Help: "Total number of HTTP requests processed",
        },
        []string{"method", "endpoint"},
    )

    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "app_request_duration_seconds",
            Help:    "Request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
)

// instrumentHandler : implements the metric instrumentation for a handler
func instrumentHandler(endpoint string, handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        timer := prometheus.NewTimer(requestDuration.WithLabelValues(r.Method, endpoint))
        defer timer.ObserveDuration()
        requestCounter.WithLabelValues(r.Method, endpoint).Inc()
        handler(w, r)
    }
}

// predict : handler function to our predict
func predict(w http.ResponseWriter, req *http.Request) {
	var pr PredictRequest
	if err := json.NewDecoder(req.Body).Decode(&pr); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()
	log.Printf("got request as : %v\n", pr)
	w.Header().Set("Content-type", "application/json")
	results := make([]float64, len(pr.Inputs))
	for i, v := range(pr.Inputs) {
		results[i] = v*float64(mrand.Int())
	}
	// TODO: v2 will be adding the sleep here for requests to take longer
	// time.Sleep(time.Second*mrand.Intn(5))
	prs := PredictResponse{
		Metadata: pr.Metadata,
		Results: results,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prs)
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
	http.Handle("/predict", instrumentHandler("predict", predict))
	log.Println("starting server at port 5000")
	if err := http.ListenAndServe(":5001", nil); err != nil {
		log.Fatalf("Server error: %s", err)
	}
}

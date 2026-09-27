# Summary

Simple go http app, exposing 2 endpoints and intrumenting prometheus metrics.

# Prometheus instrumentation

Adding the whole go info, plus two custom metrics:

* `app_http_requests_total`: Counter for the total number of requests for endpoints. 
* `app_request_duration_seconds`: Histogram metric with broken by buckets for the request duration. 

Both metrics will be used on scenarios for our `Argo-Rollouts` tests, specially the `app_request_duration_seconds`.

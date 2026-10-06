# Summary

Simple go http app, exposing 2 endpoints and intrumenting prometheus metrics.

# Prometheus instrumentation

Adding the whole go info, plus two custom metrics:

* `app_http_requests_total`: Counter for the total number of requests for endpoints. 
* `app_request_duration_seconds`: Histogram metric with broken by buckets for the request duration. 

Both metrics will be used on scenarios for our `Argo-Rollouts` tests, specially the `app_request_duration_seconds`.

Example call to the app:

```bash
curl -v http://localhost:5001/predict -d '{"metadata":{"version":"v1.0.0","caller":"caller-app"},"inputs":[1.2,1.3,4.1]}'
```

Response:

```json
{"metadata":{"version":"v1.0.0","caller":"caller-app"},"results":[0,1.3,8.2],"server_version":"v1.0.0"}
```

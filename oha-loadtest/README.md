# oha request

```bash
oha http://localhost:5001/predict \
  -m POST \
  -H "Content-Type: application/json" \
  -Z requests.ndjson \
  -q 10 \
  -z 1m \
  -c 10 \
  --no-tui
```

And we have a summary:

```bash
Summary:
  Success rate: 100.00%
  Total:        60004.8762 ms
  Slowest:      9.3436 ms
  Fastest:      0.2313 ms
  Average:      0.7224 ms
  Requests/sec: 10.0325

  Total data:   66.94 KiB
  Size/request: 114 B
  Size/sec:     1.12 KiB

Response time histogram:
  0.231 ms [1]   |
  1.143 ms [560] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  2.054 ms [32]  |■
  2.965 ms [2]   |
  3.876 ms [2]   |
  4.787 ms [2]   |
  5.699 ms [0]   |
  6.610 ms [0]   |
  7.521 ms [0]   |
  8.432 ms [0]   |
  9.344 ms [1]   |

Response time distribution:
  10.00% in 0.3520 ms
  25.00% in 0.4788 ms
  50.00% in 0.6520 ms
  75.00% in 0.8526 ms
  90.00% in 1.0465 ms
  95.00% in 1.2561 ms
  99.00% in 2.4547 ms
  99.90% in 9.3436 ms
  99.99% in 9.3436 ms


Details (average, fastest, slowest):
  DNS+dialup:   0.5502 ms, 0.3170 ms, 0.7598 ms
  DNS-lookup:   0.0548 ms, 0.0281 ms, 0.1447 ms

Status code distribution:
  [200] 600 responses

Error distribution:
  [2] aborted due to deadline
```

# oha docker and cronjob

To actually push a load test, you can just port-forward requests to your k8s service, as behind the scenes it only picks up one of the services to route your requests.

I've decided to create a docker image, and a cronjob to trigger the requests and have a constant flow of requests to my tests.

So the [docker file](./Dockerfile) here is mainly just putting `oha` on a container image, that will be used on my [cronjob](../k8s-deployments/app-argo-rollouts.yaml#95).

For each use case I'll be patching the cronjob wth certain arguments, so you can see on my [`kustomization.yaml`](../k8s-deployments/overlays/canary/kustomization.yaml#14) file that I am patching how to do the calls and to which service.

Building the image:

```bash
podman build . -t oha-loadtest:latest
```

Or build directly on your remote podman for example:

```bash
podman --connection <remote-host> build . -t oha-loadtest:latest
```

**Ps: Services are only doing round-robin calls to the underline pods / endpoints, so there's no rule by default that looks if a pod is overloaded or not.**

## References:

* <a href="https://v1-32.docs.kubernetes.io/docs/reference/kubectl/generated/kubectl_port-forward/#synopsis" target="_blank">kubectl port-forward</a>

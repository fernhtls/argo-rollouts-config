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

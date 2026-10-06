# Builder
FROM cgr.dev/chainguard/go:latest-dev AS builder
WORKDIR /app
COPY ./app/go.* ./
RUN go mod download 
COPY ./app/config.yaml ./
COPY ./app/main.go ./
RUN CGO_ENABLED=0 go build .
# Running env with chainguard
FROM cgr.dev/chainguard/static:latest
USER 1001
COPY --from=builder /app/app /app
ENTRYPOINT ["/app"]

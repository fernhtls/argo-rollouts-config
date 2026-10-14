FROM cgr.dev/chainguard/wolfi-base AS fetch
RUN apk add --no-cache curl
RUN curl -fsSL -o /oha \
      https://github.com/hatoo/oha/releases/download/v1.16.0/oha-linux-amd64 \
 && chmod +x /oha

FROM cgr.dev/chainguard/static:latest
COPY --from=fetch /oha /usr/local/bin/oha
COPY requests.ndjson .
ENTRYPOINT ["/usr/local/bin/oha"]

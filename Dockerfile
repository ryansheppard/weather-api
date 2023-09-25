FROM golang:1.21 AS builder

WORKDIR /src
COPY . .

RUN apt-get update && apt-get install -y git
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' -tags timetzdata

FROM scratch

COPY --from=builder /go/bin/weather /weather
COPY --from=golang:1.21 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


# API
EXPOSE 1323
# Prometheus
EXPOSE 1324

ENTRYPOINT ["/weather"]

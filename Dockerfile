FROM golang:alpine AS builder


COPY . .

RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/prometheus-metrics-app

FROM scratch

COPY --from=builder /go/bin/prometheus-metrics-app /go/bin/prometheus-metrics-app

EXPOSE 8080

ENTRYPOINT ["/go/bin/prometheus-metrics-app"]

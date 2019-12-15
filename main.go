package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	DEFAULTSERVERNAME = "prometheus-example-server"
)

var (
	serverName = DEFAULTSERVERNAME

	// counter type metrics
	counterLabels = map[string]string{
		"server": serverName,
	}
	counterTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "casek14",
			Subsystem:   "promexampleapp",
			Name:        "http_requests_total",
			Help:        "Total number or requests received by the server",
			ConstLabels: counterLabels,
		})
	counter200Labels = map[string]string{
		"server": serverName,
		"method": "get",
		"code":   "200",
	}
	counter200 = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "casek14",
			Subsystem:   "promexampleapp",
			Name:        "http_requests_200",
			Help:        "Total number of 200 response codes",
			ConstLabels: counter200Labels,
		})

	counter400Labels = map[string]string{
		"server": serverName,
		"method": "get",
		"code":   "400",
	}
	counter400 = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "casek14",
			Subsystem:   "promexampleapp",
			Name:        "http_requests_400",
			Help:        "Total number of 400 response codes",
			ConstLabels: counter400Labels,
		})

	counter500Labels = map[string]string{
		"server": serverName,
		"method": "get",
		"code":   "500",
	}
	counter500 = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "casek14",
			Subsystem:   "promexampleapp",
			Name:        "http_requests_500",
			Help:        "Total number of 500 response codes",
			ConstLabels: counter500Labels,
		})
)

func main() {

	rand.Seed(time.Now().Unix())

	http.Handle("/metrics", promhttp.Handler())

	prometheus.MustRegister(counterTotal)
	prometheus.MustRegister(counter200)
	prometheus.MustRegister(counter400)
	prometheus.MustRegister(counter500)

	go func() {
		for {
			counterTotal.Add(rand.Float64() * 5)
			counter200.Add(rand.Float64() * 30)
			counter400.Add(rand.Float64() * 10)
			counter500.Add(rand.Float64() *2)
			time.Sleep(time.Second * 4)
		}
	}()
log.Fatal(http.ListenAndServe(":8080",nil))
}

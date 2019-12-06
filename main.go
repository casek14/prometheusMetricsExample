package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

const (
	DEFAULTSERVERNAME = "prometheus-example-server"
)

func main() {
	p := NewPrometheusExampleServer("Karel")
	go func() {
		for {

			p.IncreaseCounters(80, 15)
			time.Sleep(time.Second)
		}
	}()
	http.HandleFunc("/metrics", p.RenderMetrics)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type PromethesExampleServer struct {
	Name                 string
	RequestsCounter      int32
	StatusCode200Counter int32
	StatusCode404Counter int32
	StatusCode500Counter int32
}

func NewPrometheusExampleServer(name string) *PromethesExampleServer {
	return &PromethesExampleServer{
		Name:                 name,
		RequestsCounter:      0,
		StatusCode200Counter: 0,
		StatusCode404Counter: 0,
		StatusCode500Counter: 0,
	}
}

func (p *PromethesExampleServer) RenderMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] Generating Metrics.\n[200-OK] %dx\n[404-NotFound] %dx\n[500-InternalServerError] %dx\n",
		time.Now().Format(time.RFC850), p.StatusCode200Counter, p.StatusCode500Counter, p.StatusCode500Counter)

	requestsCounterMetricsHelp := "#HELP casek_http_requests_total Total number or requests received by the server\n"
	requestsCounterMetricsType := "#TYPE casek_http_requests_total counter\n"
	requestsCounterMetric := fmt.Sprintf("casek_http_requests_total{server=\"%s\"} %d\n", p.Name, p.RequestsCounter)

	statusCode200CounterMetricsHelp := "#HELP casek_http_requests_total Total number of 200 response codes\n"
	statusCode200CounterMetricsType := "#TYPE casek_http_requests_total counter\n"
	statusCode200CounterMetric := fmt.Sprintf("casek_http_requests_total{method=\"get\",code=\"200\",server=\"%s\"} %d\n", p.Name, p.StatusCode200Counter)

	statusCode400CounterMetricsHelp := "#HELP casek_http_requests_total Total number of 400 response codes\n"
	statusCode400CounterMetricsType := "#TYPE casek_http_requests_total counter\n"
	statusCode400CounterMetric := fmt.Sprintf("casek_http_requests_total{method=\"get\",code=\"400\",server=\"%s\"} %d\n", p.Name, p.StatusCode404Counter)

	statusCode500CounterMetricsHelp := "#HELP casek_http_requests_total Total number of 500 response codes\n"
	statusCode500CounterMetricsType := "#TYPE casek_http_requests_total counter\n"
	statusCode500CounterMetric := fmt.Sprintf("casek_http_requests_total{method=\"get\",code=\"500\",server=\"%s\"} %d\n", p.Name, p.StatusCode500Counter)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memoryUsageHelp := "#HELP casek_memory_usage Bytes of allocated heap objects\n"
	memoryUsageType := "#TYPE casek_memory_usage gauge\n"
	memoryUsage := fmt.Sprintf("casek_memory_usage{server=\"%s\" %d}\n", p.Name, m.HeapIdle/1024)

	_, _ = fmt.Fprintf(w, requestsCounterMetricsHelp)
	_, _ = fmt.Fprintf(w, requestsCounterMetricsType)
	_, _ = fmt.Fprintf(w, requestsCounterMetric)

	_, _ = fmt.Fprintf(w, statusCode200CounterMetricsHelp)
	_, _ = fmt.Fprintf(w, statusCode200CounterMetricsType)
	_, _ = fmt.Fprintf(w, statusCode200CounterMetric)

	_, _ = fmt.Fprintf(w, statusCode400CounterMetricsHelp)
	_, _ = fmt.Fprintf(w, statusCode400CounterMetricsType)
	_, _ = fmt.Fprintf(w, statusCode400CounterMetric)

	_, _ = fmt.Fprintf(w, statusCode500CounterMetricsHelp)
	_, _ = fmt.Fprintf(w, statusCode500CounterMetricsType)
	_, _ = fmt.Fprintf(w, statusCode500CounterMetric)

	_, _ = fmt.Fprintf(w, memoryUsageHelp)
	_, _ = fmt.Fprintf(w, memoryUsageType)
	_, _ = fmt.Fprintf(w, memoryUsage)

}

// randomly increase counter on one of the counters
// stc200Max is number between 0 and 100 and represents
// % of 200 response status code
//
// stc400Max is number between 0 and (100 - stc200Max)
// represents % of 400 response status code
//
// Number of 500 status code responses will be 100 - (stc200Max + stc400Max)
func (p *PromethesExampleServer) IncreaseCounters(stc200Max int, stc400Max int) {
	if (stc200Max + stc400Max) > 100 {
		stc400Max = 100 - stc200Max - 2
		// If 200 stc + 400 stc is greater than 100,o then 400 is 100 - stc200Max -2
		// and percentage of 500 response status code is 2
	}
	// 200 status code response percentage
	switch n := getRandomNumber(); {
	case n <= stc200Max:
		p.StatusCode200Counter++
	case n > stc200Max && n <= (stc200Max+stc400Max):
		p.StatusCode404Counter++
	default:
		p.StatusCode500Counter++
	}
	p.RequestsCounter++
}

// Generate random number between 0 - 100
func getRandomNumber() int {
	rand.Seed(time.Now().UnixNano())
	number := rand.Intn(101)
	return number
}

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "request_count_total",
		Help: "Counter of HTTP requests.",
	})
	errorCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "request_error_count_total",
		Help: "Counter of HTTP requests resulting in an error.",
	})
)

func main() {
	requestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Inc()
		rand.Seed(time.Now().UnixNano())
		switch rand.Intn(3) {
		case 0:
			log.Println("OK")
			fmt.Fprint(w, "OK")
		case 1:
			log.Println("Normal")
			fmt.Fprint(w, "Normal")
		case 2:
			log.Println("Error")
			errorCount.Inc()
			fmt.Fprint(w, "Error")
		}
	})
	http.Handle("/", requestHandler)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

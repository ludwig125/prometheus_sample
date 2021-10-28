package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	executeCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "batch_count_total",
		Help: "Counter of execute.",
	})
	errorCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "batch_error_count_total",
		Help: "Counter of execute resulting in an error.",
	})
)

func main() {
	executeCount.Inc()
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

	log.Fatal(http.ListenAndServe(":8080", nil))
}

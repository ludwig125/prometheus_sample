package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	duration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "batch_duration_seconds",
		Help: "The duration of last batch in seconds.",
	})
	executeCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "batch_count_total",
		Help: "Counter of execute.",
	})
	okCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "batch_ok_count_total",
		Help: "Counter of ok execute.",
	})
	normalCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "batch_normal_count_total",
		Help: "Counter of normal execute.",
	})
	errorCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "batch_error_count_total",
		Help: "Counter of execute resulting in an error.",
	})
	missingCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "batch_missing_count_total",
		Help: "Counter of execute resulting in an missing.",
	})
)

func main() {
	var oks int
	var normals int
	var errors int

	registry := prometheus.NewRegistry()
	registry.MustRegister(duration, executeCount, okCount, normalCount, errorCount, missingCount)

	pusher := push.New("http://localhost:9091", "my_batch_job").Gatherer(registry)

	start := time.Now()

	for i := 0; i < 100; i++ {
		executeCount.Inc()
		rand.Seed(time.Now().UnixNano())
		switch rand.Intn(3) {
		case 0:
			oks++
			okCount.Inc()
		case 1:
			normals++
			normalCount.Inc()
		case 2:
			errors++
			errorCount.Inc()
		}

		time.Sleep(10 * time.Millisecond)
	}
	d := time.Since(start).Seconds()
	fmt.Println("duration:", d)
	duration.Set(d)

	fmt.Printf("ok: %d, normal: %d, error: %d\n", oks, normals, errors)
	if err := pusher.
		Push(); err != nil {
		fmt.Println(err)
	}

}

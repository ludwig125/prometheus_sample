package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
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
)

func main() {
	for i := 0; i < 100; i++ {

		executeCount.Inc()
		rand.Seed(time.Now().UnixNano())
		switch rand.Intn(3) {
		case 0:
			log.Println("OK")
			okCount.Inc()
		case 1:
			log.Println("Normal")
			normalCount.Inc()
		case 2:
			log.Println("Error")
			errorCount.Inc()
		}

		pusher := push.New("http://localhost:9091", "my_batch_job")
		if err := pusher.
			Collector(executeCount).
			Grouping("status", "all").
			Push(); err != nil {
			fmt.Println(err)
		}
		if err := pusher.
			Collector(okCount).
			Grouping("status", "ok").
			Push(); err != nil {
			fmt.Println(err)
		}
		if err := pusher.
			Collector(normalCount).
			Grouping("status", "normal").
			Push(); err != nil {
			fmt.Println(err)
		}
		if err := pusher.
			Collector(errorCount).
			Grouping("status", "error").
			Push(); err != nil {
			fmt.Println(err)
		}
	}

	// if err := push.New("http://localhost:9091", "my_batch_job").
	// 	Collector(executeCount).
	// 	Grouping("db", "customers").
	// 	Push(); err != nil {
	// 	fmt.Println(err)
	// }

	// time.Sleep(100 * time.Second)
}

https://prometheus.io/download/

```
wget https://github.com/prometheus/prometheus/releases/download/v2.31.0-rc.0/prometheus-2.31.0-rc.0.linux-amd64.tar.gz

tar -xzf prometheus-2.31.0-rc.0.linux-amd64.tar.gz

```

対象のPortに書き換える
```
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "prometheus"

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
      # - targets: ["localhost:9090"]
      - targets: ["localhost:8080"]
```

prometheusを起動

```
$./prometheus
ts=2021-10-29T21:49:48.719Z caller=main.go:406 level=info msg="No time or size retention was set so using the default time retention" duration=15d
ts=2021-10-29T21:49:48.719Z caller=main.go:444 level=info msg="Starting Prometheus" version="(version=2.31.0-rc.0, branch=HEAD, revision=21834bca6b5e44566602ea9315c8088dd82e5fad)"
ts=2021-10-29T21:49:48.720Z caller=main.go:449 level=info build_context="(go=go1.17.2, user=root@9ea31d6cef89, date=20211022-15:01:14)"
ts=2021-10-29T21:49:48.720Z caller=main.go:450 level=info host_details="(Linux 4.19.128-microsoft-standard #1 SMP Tue Jun 23 12:58:10 UTC 2020 x86_64 DESKTOP-4ND5CO6 localdomain)"
ts=2021-10-29T21:49:48.720Z caller=main.go:451 level=info fd_limits="(soft=1024, hard=4096)"
ts=2021-10-29T21:49:48.720Z caller=main.go:452 level=info vm_limits="(soft=unlimited, hard=unlimited)"
ts=2021-10-29T21:49:48.722Z caller=web.go:542 level=info component=web msg="Start listening for connections" address=0.0.0.0:9090
ts=2021-10-29T21:49:48.723Z caller=main.go:839 level=info msg="Starting TSDB ..."
ts=2021-10-29T21:49:48.724Z caller=tls_config.go:195 level=info component=web msg="TLS is disabled." http2=false
ts=2021-10-29T21:49:48.728Z caller=head.go:479 level=info component=tsdb msg="Replaying on-disk memory mappable chunks if any"
ts=2021-10-29T21:49:48.728Z caller=head.go:513 level=info component=tsdb msg="On-disk memory mappable chunks replay completed" duration=2.2µs
ts=2021-10-29T21:49:48.728Z caller=head.go:519 level=info component=tsdb msg="Replaying WAL, this may take a while"
ts=2021-10-29T21:49:48.728Z caller=head.go:590 level=info component=tsdb msg="WAL segment loaded" segment=0 maxSegment=0
ts=2021-10-29T21:49:48.728Z caller=head.go:596 level=info component=tsdb msg="WAL replay completed" checkpoint_replay_duration=20µs wal_replay_duration=386.3µs total_replay_duration=423.8µs
ts=2021-10-29T21:49:48.729Z caller=main.go:866 level=info fs_type=EXT4_SUPER_MAGIC
ts=2021-10-29T21:49:48.729Z caller=main.go:869 level=info msg="TSDB started"
ts=2021-10-29T21:49:48.729Z caller=main.go:996 level=info msg="Loading configuration file" filename=prometheus.yml
ts=2021-10-29T21:49:48.730Z caller=main.go:1033 level=info msg="Completed loading of configuration file" filename=prometheus.yml totalDuration=544µs db_storage=600ns remote_storage=1.2µs web_handler=300ns query_engine=600ns scrape=221µs scrape_sd=13.7µs notify=33.2µs notify_sd=25.1µs rules=800ns
ts=2021-10-29T21:49:48.730Z caller=main.go:811 level=info msg="Server is ready to receive web requests."
```


http://localhost:8080/ を見ると以下の通り

![image](https://user-images.githubusercontent.com/18366858/139146611-d5cac80f-5782-47c3-b14b-8c347b29f9a0.png)


http://localhost:8080/metrics

![image](https://user-images.githubusercontent.com/18366858/139337658-ed33d8d8-efc1-4eb6-93e4-67736b58cd8d.png)


# Pushgateway

https://github.com/prometheus/pushgateway

https://qiita.com/MetricFire/items/c4753396259923a0c9e2

https://kazuhira-r.hatenablog.com/entry/2019/06/02/235307



## install


https://github.com/prometheus/pushgateway/releases

```
wget https://github.com/prometheus/pushgateway/releases/download/v1.4.2/pushgateway-1.4.2.linux-amd64.tar.gz

tar -xzf pushgateway-1.4.2.linux-amd64.tar.gz

rm pushgateway-1.4.2.linux-amd64.tar.gz
```

```
$./pushgateway
ts=2021-10-29T21:50:55.756Z caller=level.go:63 level=info msg="starting pushgateway" version="(version=1.4.2, branch=HEAD, revision=99981d7be923ab18d45873e9eaa3d2c77477b1ef)"
ts=2021-10-29T21:50:55.759Z caller=level.go:63 level=info build_context="(go=go1.16.9, user=root@f68dbd4cbcde, date=20211011-17:51:55)"
ts=2021-10-29T21:50:55.760Z caller=level.go:63 level=info listen_address=:9091
ts=2021-10-29T21:50:55.761Z caller=level.go:63 level=info msg="TLS is disabled." http2=false
```



## reference

https://kazuhira-r.hatenablog.com/entry/2019/06/02/235307
https://qiita.com/MetricFire/items/c4753396259923a0c9e2
https://kobatako.hatenablog.com/entry/2020/01/07/231108
https://it-engineer.hateblo.jp/entry/2019/01/12/105700
https://stackoverflow.com/questions/37611754/how-to-push-metrics-to-prometheus-using-client-golang
https://www.robustperception.io/choosing-your-pushgateway-grouping-key



```
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

```

結果
pushgateway
- http://localhost:9091/metrics

```
# HELP batch_count_total Counter of execute.
# TYPE batch_count_total counter
batch_count_total{instance="",job="my_batch_job",status="all"} 100
batch_count_total{instance="",job="my_batch_job",status="error"} 100
batch_count_total{instance="",job="my_batch_job",status="normal"} 100
batch_count_total{instance="",job="my_batch_job",status="ok"} 100
# HELP batch_error_count_total Counter of execute resulting in an error.
# TYPE batch_error_count_total counter
batch_error_count_total{instance="",job="my_batch_job",status="error"} 31
# HELP batch_normal_count_total Counter of normal execute.
# TYPE batch_normal_count_total counter
batch_normal_count_total{instance="",job="my_batch_job",status="error"} 36
batch_normal_count_total{instance="",job="my_batch_job",status="normal"} 36
# HELP batch_ok_count_total Counter of ok execute.
# TYPE batch_ok_count_total counter
batch_ok_count_total{instance="",job="my_batch_job",status="error"} 33
batch_ok_count_total{instance="",job="my_batch_job",status="normal"} 33
batch_ok_count_total{instance="",job="my_batch_job",status="ok"} 33
```

http://localhost:9091/#

![image](https://user-images.githubusercontent.com/18366858/139559832-4dc5b3d7-a7e7-43cf-9ef4-f806c48a98cc.png)


![image](https://user-images.githubusercontent.com/18366858/139559850-c5c829a5-abf3-4b75-83e5-8bab97b6439a.png)

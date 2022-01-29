## server の場合

以下のような Prometheus 用のメトリクスを 8080/metrics に出力するサーバを想定する

```go
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

```

起動してみる

```
cd server

go run server.go
```

http://localhost:8080/ にリクエストするたびに結果が変わるのが分かる

![image](https://user-images.githubusercontent.com/18366858/151625269-2acda0b3-9f0b-4a7a-b604-33ad49e04730.png)
![image](https://user-images.githubusercontent.com/18366858/151625311-5fd85488-0bda-4c62-9afb-ed18167ff529.png)

このメトリクスは以下に表示される

http://localhost:8080/metrics

![image](https://user-images.githubusercontent.com/18366858/139337658-ed33d8d8-efc1-4eb6-93e4-67736b58cd8d.png)

## prometheus のインストールと起動

https://prometheus.io/download/

```
wget https://github.com/prometheus/prometheus/releases/download/v2.31.0-rc.0/prometheus-2.31.0-rc.0.linux-amd64.tar.gz

tar -xzf prometheus-2.31.0-rc.0.linux-amd64.tar.gz

```

上で起動した server の Port に監視対象を書き換える

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

prometheus を起動

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

http://localhost:9090/ を見ると以下の通り Prometheus の UI が見られる

![image](https://user-images.githubusercontent.com/18366858/139146611-d5cac80f-5782-47c3-b14b-8c347b29f9a0.png)

server 側で出したメトリクスをグラフ化することもできる

![image](https://user-images.githubusercontent.com/18366858/151625755-26ff9409-6a50-43f3-8f7c-86557e314d97.png)

以上がサーバの場合

# Pushgateway

batch の場合, Prometheus が見るころには Batch 処理は終わっていてメトリクスが取れない問題が起きる

つまりさっさと終わってしまった Batch の代わりに`誰か`が Batch 処理が出したメトリクスを拾って取っておかないと、Prometheus に見てもらえない

ということでその誰かが Pushgateway にあたる

## Pushgateway

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

Batch の場合、pushgateway を先に起動しておく必要がある

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

# Batch の例

先に prometheus と pushgateway を起動した状態で

以下の Batch を起動してみる

## Add の例

https://pkg.go.dev/github.com/prometheus/client_golang/prometheus/push#Pusher.Add

をそのまま実行してみる(push 先は localhost に変える)

```go
package main

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	completionTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_last_completion_timestamp_seconds",
		Help: "The timestamp of the last completion of a DB backup, successful or not.",
	})
	successTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_last_success_timestamp_seconds",
		Help: "The timestamp of the last successful completion of a DB backup.",
	})
	duration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_duration_seconds",
		Help: "The duration of the last DB backup in seconds.",
	})
	records = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_records_processed",
		Help: "The number of records processed in the last DB backup.",
	})
)

func performBackup() (int, error) {
	// Perform the backup and return the number of backed up records and any
	// applicable error.
	// ...
	return 42, nil
}

func main() {
	// We use a registry here to benefit from the consistency checks that
	// happen during registration.
	registry := prometheus.NewRegistry()
	registry.MustRegister(completionTime, duration, records)
	// Note that successTime is not registered.

	pusher := push.New("http://localhost:9091", "db_backup").Gatherer(registry)

	start := time.Now()
	n, err := performBackup()
	records.Set(float64(n))
	// Note that time.Since only uses a monotonic clock in Go1.9+.
	duration.Set(time.Since(start).Seconds())
	completionTime.SetToCurrentTime()
	if err != nil {
		fmt.Println("DB backup failed:", err)
	} else {
		// Add successTime to pusher only in case of success.
		// We could as well register it with the registry.
		// This example, however, demonstrates that you can
		// mix Gatherers and Collectors when handling a Pusher.
		pusher.Collector(successTime)
		successTime.SetToCurrentTime()
	}
	// Add is used here rather than Push to not delete a previously pushed
	// success timestamp in case of a failure of this backup.
	if err := pusher.Add(); err != nil {
		fmt.Println("Could not push to Pushgateway:", err)
	}
}
```

```
[~/go/src/github.com/ludwig125/prometheus_sample] $go run batch_add/add.go
```

http://localhost:9091/ Pushgateway を見ると以下のようになる

![image](https://user-images.githubusercontent.com/18366858/151678769-7f405f8e-c1d1-4b56-b662-047a366fa7a0.png)

Prometheus は以下のようにそれぞれメトリクスが取得できている

http://localhost:9090/graph?g0.expr=db_backup_last_completion_timestamp_seconds&g0.tab=1&g0.stacked=0&g0.show_exemplars=0&g0.range_input=1h&g1.expr=db_backup_last_success_timestamp_seconds&g1.tab=1&g1.stacked=0&g1.show_exemplars=0&g1.range_input=1h&g2.expr=db_backup_duration_seconds&g2.tab=1&g2.stacked=0&g2.show_exemplars=0&g2.range_input=1h&g3.expr=db_backup_records_processed&g3.tab=1&g3.stacked=0&g3.show_exemplars=0&g3.range_input=1h

![image](https://user-images.githubusercontent.com/18366858/151678801-cf403bf2-afa1-4113-9f46-1f4d975f8fdc.png)

## Push

次に Push を見てみる
https://pkg.go.dev/github.com/prometheus/client_golang/prometheus/push#Pusher.Push

をそのまま実行する

```go
package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

func main() {
	completionTime := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_last_completion_timestamp_seconds",
		Help: "The timestamp of the last successful completion of a DB backup.",
	})
	completionTime.SetToCurrentTime()
	if err := push.New("http://localhost:9091", "db_backup").
		Collector(completionTime).
		Grouping("db", "customers").
		Push(); err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	}
}

```

```
[~/go/src/github.com/ludwig125/prometheus_sample] $go run batch_push/push.go
```

pushgateway と prometheus は以下の通り

![image](https://user-images.githubusercontent.com/18366858/151679145-2ec9653b-616c-4379-b4b0-99584d34bb74.png)
![image](https://user-images.githubusercontent.com/18366858/151679163-6ff34b4e-03c7-49dc-9b62-c473f723509b.png)

Grouping の役割を理解するために
上で作った、`batch_push/push.go`を少し書き換えて、
`Grouping("db", "customers")`の部分を`Grouping("db", "producers")`とした Batch を別に実行させてみる

```
[~/go/src/github.com/ludwig125/prometheus_sample] $go run batch_push2/push.go
```

以下のように`customers`と`producers`が分かれて登録されるので、異なる種類の DB の結果を別に扱うのに便利そう

![image](https://user-images.githubusercontent.com/18366858/151679323-ba404fc1-d1d9-4084-8cb8-ff96e693f58d.png)
![image](https://user-images.githubusercontent.com/18366858/151679331-3756557c-4f5b-4f1f-988e-5da222bd26fc.png)

### Add と Push の違い

https://pkg.go.dev/github.com/prometheus/client_golang/prometheus/push#Pusher.Add

```
Addはpushと同じように動作しますが、同じ名前（および同じジョブや他のグループ化ラベル）を持つ以前にpushされたメトリクスのみが置き換えられます。(PushgatewayへのプッシュにはHTTPメソッド "POST "を使用します)。
```

https://pkg.go.dev/github.com/prometheus/client_golang/prometheus/push#Pusher.Push

```
Pushは、このPusherに追加されたすべてのCollectorとGatherersからすべてのメトリクスを収集/集計します。次に、設定されたジョブ名と追加されたグループ化ラベルをグループ化キーとして、この Pusher の作成時に設定された Pushgateway にそれらをプッシュします。同じジョブや他のグルーピング・ラベルを持つ、以前にプッシュされたすべてのメトリクスは、この呼び出しによってプッシュされたメトリクスに置き換えられます。(PushgatewayへのプッシュにはHTTPメソッド "PUT "を使用します)。
```

ということで、

- Add => POST: リソースの作成
- Push => PUT: リソースの作成、リソースの置換

と理解。細かい動作確認はしていない

## 自分なりの Batch

以下のような Batch を考える

ここでは、先ほどの Server のリクエストにあたるものを 100 回実行させるバッチを考える
処理の中で Pushgateway の 9091 宛にメトリクスを飛ばしている

```go
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
)

func main() {
	var oks int
	var normals int
	var errors int

	registry := prometheus.NewRegistry()
	registry.MustRegister(duration, executeCount)

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
		Collector(okCount).
		Collector(normalCount).
		Collector(errorCount).
		Push(); err != nil {
		fmt.Println(err)
	}

}

```

前述の Pushgateway のサンプルを参考に、
`batch_duration_seconds`や`batch_count_total`という、絶対にメトリクスとして取得されるものは
以下のように、事前に resistry として登録しておく

```go
registry := prometheus.NewRegistry()
registry.MustRegister(duration, executeCount)

pusher := push.New("http://localhost:9091", "my_batch_job").Gatherer(registry)
```

ok 数や error 数のように、あとから登録してもいいものは後から`Collector`を使って登録した

registry を使わずに、以下のように全部後から登録してもいいのではと思ったけどどうだろう

```go
pusher := push.New("http://localhost:9091", "my_batch_job")
if err := pusher.
	Collector(duration).
	Collector(executeCount).
	Collector(okCount).
	Collector(normalCount).
	Collector(errorCount).
	Push(); err != nil {
	fmt.Println(err)
}
```

結果

```
[~/go/src/github.com/ludwig125/prometheus_sample] $go run batch_mysample/mysample.go
duration: 1.0421209
ok: 34, normal: 26, error: 40
[~/go/src/github.com/ludwig125/prometheus_sample] $
```

pushgateway 側のメトリクスに Server の場合の時のメトリクスが表示されている

- http://localhost:9091/metrics

```
# HELP batch_count_total Counter of execute.
# TYPE batch_count_total counter
batch_count_total{instance="",job="my_batch_job"} 100
# HELP batch_duration_seconds The duration of last batch in seconds.
# TYPE batch_duration_seconds gauge
batch_duration_seconds{instance="",job="my_batch_job"} 1.0421209
# HELP batch_error_count_total Counter of execute resulting in an error.
# TYPE batch_error_count_total counter
batch_error_count_total{instance="",job="my_batch_job"} 40
# HELP batch_normal_count_total Counter of normal execute.
# TYPE batch_normal_count_total counter
batch_normal_count_total{instance="",job="my_batch_job"} 26
# HELP batch_ok_count_total Counter of ok execute.
# TYPE batch_ok_count_total counter
batch_ok_count_total{instance="",job="my_batch_job"} 34
```

http://localhost:9091 の Pushgateway の UI は以下のようになる

![image](https://user-images.githubusercontent.com/18366858/151680103-9fe8d3f1-75f9-4593-8fa8-338b02b434ad.png)
![image](https://user-images.githubusercontent.com/18366858/151680114-1e2697e7-95cc-4524-8552-df1eab828d7d.png)

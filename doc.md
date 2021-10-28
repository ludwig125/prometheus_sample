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
```

## 为 HTTPServer 添加 0-2 秒的随机延时

  ```go
  delay := randInt(0, 2000)
  time.Sleep(time.Microsecond * time.Duration(delay))
  ```

##  为 HTTPServer 项目添加延时 Metric：httpserver/metrics 

  ```go
  func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
  	return prometheus.NewHistogramVec(
  		prometheus.HistogramOpts{
  			Namespace: namespace,
  			Name:      "execution_latency_seconds",
  			Help:      help,
  			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
  		}, []string{"step"},
  	)
  }
  ```

## 将 HTTPServer 部署至测试集群，并完成 Prometheus 配置；

  ```bash
  helm repo add grafana https://grafana.github.io/helm-charts
  helm repo update
  kubectl create ns monitoring  # 创建监控命名空间
  helm install loki-stack -n monitoring --set grafana.enabled=true --set prometheus.enabled=true grafana/loki-stack
  ```

## 从 Promethus 界面中查询延时指标数据：

  ```bash
  kubectl port-forward -n monitoring svc/loki-stack-prometheus-server 8000:80 # 本地端口转发，使用浏览器访问 http://127.0.0.1:8000
  ```

## 创建一个 Grafana Dashboard 展现延时分配情况：

  ```bash
  kubectl get secret -n monitoring loki-stack-grafana -ojsonpath={.data.admin-password}|base64 -d
  ```
### 1. 安装 istio
```shell
curl -L https://istio.io/downloadIstio | sh -
```
进入istio目录并把 istioctl 复制到 /usr/bin 目录
```shell
cd istio-1.16.1

cp bin/istioctl /usr/local/bin
```

安装 istio 到 k8s 集群
```shell
istioctl install --set profile=demo -y
```

安装 istio 后，k8s中增加的 istio 的 namespace、pod、service
```shell
#  查看 istio 的 namespace istio-system
kubectl get ns
# 查看 istio 的 pod 和 service
kubectl get po,svc -n istio-system
```

### 2. 安装 jaeger，设置 tracing 采样比为100%。

```shell
kubectl apply -f jaeger.yaml

kubectl get pod -n istio-system

kubectl edit configmap istio -n istio-system 
```

### 3.改造 httpserver
区分service1、service2、service3，并为 service1 和 service2 增加调用下游服务与请求头向下游传递逻辑。

### 4.构建 service1、2、3 镜像
```shell
docker build -f service1/Dockerfile -t service1:1.0.0 service1

docker build -f service2/Dockerfile -t service2:1.0.0 service2

docker build -f service3/Dockerfile -t service3:1.0.0 service3
```

### 5.创建带有istio注入的namespace、service1、2、3 的 deployment、service
```shell
kubectl apply -f httpserver.yaml 
kubectl get po,svc -n httpserver
```

### 6.创建证书
```shell
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=fcx Inc./CN=*.fcx.org' -keyout fcx.org.
```

### 7.在 istio 中配置证书 secret
```shell
kubectl create -n istio-system secret tls fcx-crt --key=fcx.org.key --cert=fcx.org.crt

kubectl get -n istio-system secret
```

### 8.配置 istio VirtualService 和 Gateway
```shell
kubectl apply -f istio-vs-gw.yaml -n httpserver
```

### 9.通过 istio Gateway 访问 httpserver
```shell
kubectl get svc -nistio-system

curl --resolve service.fcx.org:443:10.98.20.109 https://service.fcx.org/service0 -v -k
```

### 10.从 jaeger dashboard 查看链路跟踪信息
```shell
istioctl dashboard jaeger --address ${IP} # `--address` 指定dashboard访问ip到虚拟机ip上，方便通宿主机访问
```


## 模块八作业第二部分
### 题目：
除了将 httpServer 应用优雅的运行在 Kubernetes 之上，我们还应该考虑如何将服务发布给对内和对外的调用方。 来尝试用 Service, Ingress 将你的服务发布给集群外部的调用方吧。 在第一部分的基础上提供更加完备的部署 spec，包括（不限于）：

Service Ingress 可以考虑的细节

如何确保整个应用的高可用。 如何通过证书保证 httpServer 的通讯安全。

### 步骤
1. 基于 part 1 作业，在 httpserver 的 namespace 部署两个 httpserver
```shell
kubectl create ns httpserver

kubectl create -f config.yaml -n httpserver

kubectl get configmap -n httpserver

kubectl create -f deploy.yaml -n httpserver

kubectl get pod -n httpserver
```
2. 配置Service: service.yaml
```shell
kubectl create -f service.yaml -n httpserver

kubectl get svc -n httpserver
```
3. 部署ingress-nginx-controller，并生成证书
4. 配置Ingress: gateway-ingress.yaml
```shell
kubectl create -f gateway-ingress.yaml -n httpserver

kubectl get ingress -n httpserver

kubectl get svc -n ingress-nginx
```
# go-http-seerver
## 启动服务
```shell
go run main.go
```
## 模块三作业
- 构件本地镜像，使用练httpserver
```shell
docker build -t yys/httpserver:v0 .  
```
![构件本地镜像](./images/Screen%20Shot%202022-10-16%20at%2022.05.01.png)

- 通过 docker 命令本地启动 httpserver
```shell
docker run -p 80:80 yys/httpserver:v0
```

- 通过 nsenter 进入容器查看 IP 配置
  * 进入容器内部  
  ```shell
  # docker exec -it <CONTAINER ID> /bin/sh 
  # 在本地报错，使用 sudo 解决
  sudo docker exec <CONTAINER ID> /bin/sh
  ```
  报错解决: https://stackoverflow.com/questions/48001082/oci-runtime-exec-failed-exec-failed-executable-file-not-found-in-path

  * 查看容器内进程 pid
  ```shell
  # pid=`lsns -t net -n|grep httpserver|awk '{print $4}'` lsns没有安装，所以用不了

  # 查看进程 PID
  ps
  ```
  * 查看 httpserver IP 配置
  ```shell
  # nsenter -t $pid -n ip addr 报错
  nsenter -t <PID>
  ip addr
  ```
![通过 nsenter 进入容器查看 IP 配置](./images/Screen%20Shot%202022-10-16%20at%2022.23.53.png)
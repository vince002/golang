

# 模块八

# 第一部分： 编写 Kubernetes 部署脚本将 httpserver 部署到 kubernetes 集群

## 思考的维度

## 1、优雅启动
readinessProbe
## 2、优雅终止
terminationGracePeriodSeconds 设置为60
复制101 中httpserver sigterm相关代码到main.go
signal.Notify 
修改Dockerfile 增加tini

## 3、资源需求和 QoS 保证
Burstable 
## 4、探活
livenessProbe
## 5、日常运维需求，日志等级
通过环境变量v设置日志等级

## 6、配置和代码分离
httpserver.yaml中环境变量例子
```
      env:
        - name: VERSION
          value: 1.0.1
```
### 7、 执行步骤
#### 1、GoLand 工程目录执行
```
make push
```
#### 2、ssh登录本地虚拟机拉取镜像
```
sudo docker login hub.docker.com --username vinceleung --password ***
sudo docker pull vinceleung/httpserver:1.0.1
```
访问不到hub.docker.com,通过离线包方式导入镜像到虚拟机
```
docker.io 域名可以pull
sudo docker pull docker.io/vinceleung/httpserver:1.0.2
 ```

```
docker save vinceleung/httpserver:1.0.1  -o httpserver101.zip

sudo docker load -i httpserver101.zip
```
#### 3、执行创建Pod
上传httpserver.yaml 文件到虚拟机/opt/app目录
执行以下命令创建Pod
```
kubectl create -f httpserver.yaml
```

#### 4、查询服务状态以及访问
```
k describe po  httpserver

curl 192.168.119.123/healthz

curl 192.168.119.123/?user=vince

```

# 第二部分： 如何将服务发布给对内和对外的调用方。
用 Service, Ingress 将的服务发布给集群外部的调用方
### 可以考虑的细节：
如何确保整个应用的高可用 </br>
如何通过证书保证 httpServer 的通讯安全

## 创建httpserver-deploy
```
k apply -f httpserver-deploy.yaml
```

## 创建httpserver-service
```
k apply -f httpserver-service.yaml
```

## 创建httpserver-gateway
```
k describe ingress httpserver-gateway
```

## 测试返回结果
```
curl -H "Host: vinceleung.com" http://192.168.119.118 -v -k
```

## 通过证书保证 httpServer 的通讯安全
```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=vinceleung.com/O=vinceleung"

kubectl create secret tls vinceleung-tls --key="tls.key" --cert="tls.crt"

k get secret vinceleung-tls -o yaml
```

# 复制yaml 到httpserver-secret.yaml
```
k delete secrets vinceleung-tls

k apply -f httpserver-secret.yaml
```

## 测试返回结果
```
curl -H "Host: vinceleung.com" https://192.168.119.121 -v -k

```









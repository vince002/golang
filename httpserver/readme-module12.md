# 模块十二
## 把httpserver服务以 Istio Ingress Gateway 的形式发布

### 1、实现安全保证

#### deploy httpserver

```
cd httpserver/specs/istio
kubectl create ns httpserver-istio
kubectl label ns httpserver-istio istio-injection=enabled
kubectl create -f httpserver.yaml -n httpserver-istio
```
```
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=vinceleung Inc./CN=*.vinceleung.io' -keyout vinceleung.io.key -out vinceleung.io.crt
kubectl create -n istio-system secret tls vinceleung-credential --key=vinceleung.io.key --cert=vinceleung.io.crt
kubectl apply -f istio-specs.yaml -n httpserver-istio
```

### check ingress ip
```
k get svc -nistio-system
istio-ingressgateway   LoadBalancer   $INGRESS_IP
```
### access the httpserver via ingress
```
curl --resolve httpsserver.vinceleung.io:443:$INGRESS_IP https://httpsserver.vinceleung.io/healthz -v -k
```


### 2、七层路由规则
httpserver/specs/istio/istio-specs.yaml 文件增加url匹配规则
```
- match:
  - uri:
  exact: "/vinceleung/"
  rewrite:
  uri: "/"
  route:
  - destination:
  host: httpserver.httpserver-istio.svc.cluster.local
  port:
  number: 80
  - 
```
```
kubectl apply -f istio-specs.yaml -n httpserver-istio
```
```
curl --resolve httpsserver.vinceleung.io:443:10.99.239.22 https://httpsserver.vinceleung.io/vinceleng?user=testI7MathchByVinceLeung -v -k
```

### 3、open tracing 的接入

应用层面处理header
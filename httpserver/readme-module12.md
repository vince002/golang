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

#### main.go deal with header

```
req, err := http.NewRequest("GET", "http://service1", nil)
if err != nil {
    fmt.Printf("%s", err)
}
lowerCaseHeader := make(http.Header)
for key, value := range r.Header {
    lowerCaseHeader[strings.ToLower(key)] = value
}
glog.Info("headers:", lowerCaseHeader)
req.Header = lowerCaseHeader
```

```
kubectl apply -f jaeger.yaml
```
update sampling rate is 100%
```
kubectl edit configmap istio -n istio-system
set tracing.sampling=100
```


#### deploy tracing
```
kubectl apply ns httpserver-istio
kubectl label ns httpserver-istio istio-injection=enabled
kubectl -n httpserver-istio apply -f httpserver.yaml
kubectl -n httpserver-istio apply -f service1.yaml
kubectl -n httpserver-istio apply -f service2.yaml
kubectl apply -f istio-specs.yaml -n httpserver-istio
```
#### check ingress ip
```
k get svc -nistio-system
istio-ingressgateway   LoadBalancer   $INGRESS_IP
```
#### access the tracing via ingress for 100 times(sampling rate is 100%)
```

curl --resolve httpsserver.vinceleung.io:443:10.99.239.22 https://httpsserver.vinceleung.io/vinceleng?user=testI7MathchByVinceLeung -v -k

```
#### check tracing dashboard

```
istioctl dashboard jaeger
```
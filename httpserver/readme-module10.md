

# 模块十

## 思考的维度
### 1、为 HTTPServer 添加 0-2 秒的随机延时
```
    main.go rootHandler 
    delay := randInt(10,2000)
	time.Sleep(time.Millisecond*time.Duration(delay))
```

### 2、为 HTTPServer 项目添加延时 Metric
```
    //注册metrics
    mux.Handle("/metrics", promhttp.Handler())
```
```
    //注册指标
    metrics/metrics.go 
    main.go main() 
    metrics.Register()
```
```
    //输出指标
    main.go rootHandler
    timer := metrics.NewTimer()
	defer timer.ObserveTotal()
```


###  3、将 HTTPServer 部署至测试集群，并完成 Prometheus 配置
```
    //specs目录的deployment.yaml
    kubectl create -f  deployment.yaml
```

访问Pod根路径

```
k get po -owide

curl 192.168.119.99
```
查看日志
```
k logs httpserver-6cbd889448-k5ssh
```
![Pod日志](images/httpserver-pod-log.png)

访问/metrics

```
curl 192.168.119.99/metrics
```
![访问Pod metrics路径](images/httpserver-pod-metrics.png)



安装loki-stack
```
helm pull grafana/loki-stack
tar -xvf loki-stack-2.4.1.tgz
cd loki-stack
```
下载安装包失败，虚拟机github访问不到
在mac下载、解压、替换文件内容复制到虚拟机

[loki-stack2.4.1下载地址](https://objects.githubusercontent.com/github-production-release-asset-2e65be/289034413/22b0f100-b6fb-11eb-9e92-293091ef64e0?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20211213%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20211213T095207Z&X-Amz-Expires=300&X-Amz-Signature=)

用Visual Studio Code替换版本
```
replace all `rbac.authorization.k8s.io/v1beta1` with `rbac.authorization.k8s.io/v1`
```

复制虚拟机器文件夹loki-stack
```
mv /opt/app/101/loki-stack /home/cadmin/.cache/helm/repository/loki-stack
```

install loki locally
```
helm upgrade --install loki ./loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
```

###  4、从 Promethus 界面中查询延时指标数据
change the loki-prometheus-server service to NodePort type and access it

[prometheus-server查看指标路径](http://192.168.34.2:31185/graph?g0.range_input=1h&g0.expr=httpserver_execution_latency_seconds_bucket&g0.tab=1)


![httpserver_execution_latency_seconds_bucket](images/prometheus_latency_seconds_bucket.png)
###  5、（可选）创建一个 Grafana Dashboard 展现延时分配情况

To get the admin password for the Grafana pod, run the following command

```
k edit svc loki-grafana

kubectl get secret loki-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

```
http://192.168.34.2:31214/
admin zYbeKgVB4yWJWGdsXfeoutiLxeV1XDzM5jUtNTOw
```

#### [Import Dashboard](http://192.168.34.2:31214/dashboard/import)

![导入dashboard-httpserver-latency](images/dashboard-httpserver-latency-import.png)

#### Create Dashboard
创建dashboard后，复制grafana-dashboard/httpserver-latency.json内容设置到JSON Model
另存为http-server-latency

### 预览dashboards-http-server-latency数据

访问httpserver根路径
```
curl 192.168.119.99
```
General /Home 界面 Dashboards 面板点击http-server-latency预览数据

根据prometheus的config配置同步指标到grafana

![alt dashboard-httpserver-latency](images/dashboard-httpserver-latency.png)









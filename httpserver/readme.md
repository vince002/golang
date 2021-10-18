

模块三作业

1、构建本地镜像
cd  /Users/coding/go/src/github.com/vince002/golang/httpserver
make release

2、推送镜像

docker login
输入用户和密码

docker push vinceleung/httpserver:1.0.0
镜像仓库地址
https://hub.docker.com/repository/docker/vinceleung/httpserver
3、通过 Docker 命令本地启动 httpserver

docker run -p 8081:80 vinceleung/httpserver:1.0.0

4、通过 nsenter 进入容器查看 IP 配置

查看容器ID 
docker ps 

查看该容器的 Pid
docker inspect -f {{.State.Pid}} 8b00567c0ed0

nsenter 命令进入该容器的网络命令空间
nsenter -n -t 3786






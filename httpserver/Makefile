export tag=1.0.0
root:
	export ROOT=github.com/vince002/golang

build:
	echo "building httpserver binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building httpserver container"
	docker build -t vinceleung/httpserver:${tag} .

push: release
	echo "pushing vinceleung/httpserver"
	docker push vinceleung/httpserver:${tag}
pull:
	sudo docker pull vinceleung/httpserver:${tag}
run:
	sudo docker run vinceleung/httpserver:${tag}

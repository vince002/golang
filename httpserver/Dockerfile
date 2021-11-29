FROM ubuntu
ENV MY_SERVICE_PORT=80
LABEL multi.label1="go" multi.label2="httpserver" other="文斯"

# Add Tini
# see https://cloud-atlas.readthedocs.io/zh_CN/latest/docker/init/docker_tini.html
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

# Add httpserver
ADD bin/amd64/httpserver /httpserver
EXPOSE 80
ENTRYPOINT ["/tini", "--"]

# Run your program under Tini
# CMD ["/your/program", "-and", "-its", "arguments"]
CMD ["/httpserver"]

# ENTRYPOINT 方式
#ENTRYPOINT /httpserver
# requiring Docker 17.05 or higher on the daemon and client
# see https://docs.docker.com/develop/develop-images/multistage-build/
# BUILD COMMAND :
# docker --build-arg RELEASE_VERSION=0.0.1 -t CloudMonitor/hawkeye:v0.0.1 .

# author: liyan
# mail: huntsman_ly@sina.com

# build frontend

FROM hub.docker.com/base/golang:1.11 as backend

COPY . /go/src/github.com/huntsman-li/resp/

RUN cd /go/src/github.com/huntsman-li/resp/ && go build .


# build release image

FROM hub.docker.com/alpine/alpine as alpine

COPY --from=backend /go/src/github.com/huntsman-li/resp/resp /backend

WORKDIR /

CMD ["./backend"]
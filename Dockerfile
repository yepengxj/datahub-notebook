# This is a local-build docker image for p2p-dl test
FROM golang:1.5
MAINTAINER Zonesan <chaizs@asiainfo.com>
EXPOSE 35800
ENV GOPATH /home/go
WORKDIR $GOPATH/src/github.com/asiainfoLDP/datahub
ADD . $GOPATH/src/github.com/asiainfoLDP/datahub
RUN go get github.com/tools/godep && \
    $GOPATH/bin/godep restore && \
    $GOPATH/bin/godep go install && \
    mv daemon /bin

entrypoint daemon --daemon

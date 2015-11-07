# This is a local-build docker image for p2p-dl test

FROM golang:1.5
MAINTAINER Zonesan <chaizs@asiainfo.com>
RUN go get github.com/asiainfoLDP/datahub

EXPOSE 35800
ENTRYPOINT $GOPATH/bin/datahub --daemon

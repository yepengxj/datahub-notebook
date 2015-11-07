# This is a local-build docker image for p2p-dl test

FROM golang:1.5
MAINTAINER Zonesan <chaizs@asiainfo.com>

RUN mkdir /home/datahub
WORKDIR /home/datahub
RUN git clone https://github.com/asiainfoLDP/datahub .
RUN go get github.com/tools/godep
RUN $GOPATH/bin/godep restore
RUN $GOPATH/bin/godep go install


EXPOSE 35800
ENTRYPOINT $GOPATH/bin/datahub --daemon

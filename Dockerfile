# This is a local-build docker image for p2p-dl test
FROM golang:1.5
MAINTAINER Zonesan <chaizs@asiainfo.com>
EXPOSE 35800
ENV GOPATH /home/go
WORKDIR $GOPATH/src/github.com/asiainfoLDP/datahub
ADD . $GOPATH/src/github.com/asiainfoLDP/datahub
RUN go get github.com/tools/godep && \
    godep restore && \
    godep go install
RUN cd $GOPATH/src/github.com/asiainfoLDP/datahub && \
    curl -s https://raw.githubusercontent.com/pote/gpm/v1.3.2/bin/gpm | bash && \
    go build && \
    mv datahub /bin

entrypoint daemon --daemon

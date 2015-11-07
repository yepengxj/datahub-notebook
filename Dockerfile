# This is a local-build docker image for p2p-dl test
FROM golang:1.5
MAINTAINER Zonesan <chaizs@asiainfo.com>
EXPOSE 35800
WORKDIR /home/go
ENV GOPATH /home/go
ADD . $GOPATH/src/github.com/asiainfoLDP/datahub
RUN cd $GOPATH/src/github.com/asiainfoLDP/datahub && \
    curl -s https://raw.githubusercontent.com/pote/gpm/v1.3.2/bin/gpm | bash && \
    go build && \
    mv datahub /bin

entrypoint daemon --daemon

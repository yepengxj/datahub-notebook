# This is a local-build docker image for p2p-dl test
FROM golang:1.5
MAINTAINER Zonesan <chaizs@asiainfo.com>
EXPOSE 59090
WORKDIR /home
RUN echo "This is a local-build docker image for test" > readme
RUN git clone https://github.com/asiainfoLDP/datahub-client
WORKDIR datahub-client
RUN go build

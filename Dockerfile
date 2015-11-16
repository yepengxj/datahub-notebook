# This is a local-build docker image for p2p-dl test

FROM golang:1.5
MAINTAINER Zonesan <chaizs@asiainfo.com>

ENV SRCPATH $GOPATH/src/github.com/asiainfoLDP/datahub 
ENV PATH $PATH:$GOPATH/bin:$SRCPATH
RUN mkdir $SRCPATH -p
WORKDIR $SRCPATH

ADD . $SRCPATH
#COPY . /home/datahub
#RUN git clone https://github.com/asiainfoLDP/datahub .
#RUN go get github.com/tools/godep
#RUN $GOPATH/bin/godep restore
#RUN $GOPATH/bin/godep go install
run mkdir /var/lib/datahub
#run tar zxvf test/repos.tar.gz -C /var/lib/datahub
run curl -s https://raw.githubusercontent.com/pote/gpm/v1.3.2/bin/gpm | bash && \
    go build

EXPOSE 35800

CMD $SRCPATH/start.sh



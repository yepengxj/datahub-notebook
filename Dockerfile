FROM jupyter/all-spark-notebook
MAINTAINER Zonesan <chaizs@asiainfo.com>


##install golang runtime
RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make 

ENV GOLANG_VERSION 1.5.1
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA1 46eecd290d8803887dec718c691cc243f2175fe0

RUN wget https://s3.cn-north-1.amazonaws.com.cn/asiainfoldp-file-backup/golang.tar.gz -O /usr/local/golang.tar.gz \
	&& tar -C /usr/local -xzf /usr/local/golang.tar.gz \
	&& rm /usr/local/golang.tar.gz \
        && ls /usr/local/go/bin/


ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

COPY golang1.5/go-wrapper /usr/local/bin/


##install datahub deamon
ENV SRCPATH $GOPATH/src/github.com/asiainfoLDP/datahub 
ENV PATH $PATH:$GOPATH/bin:$SRCPATH/bin
RUN mkdir $SRCPATH -p
WORKDIR $SRCPATH

ADD ./datahub_d/ $SRCPATH
run ls $SRCPATH
#COPY . /home/datahub
#RUN git clone https://github.com/asiainfoLDP/datahub .
#RUN go get github.com/tools/godep
#RUN $GOPATH/bin/godep restore
#RUN $GOPATH/bin/godep go install
run mkdir /var/lib/datahub
#run tar zxvf test/repos.tar.gz -C /var/lib/datahub
RUN curl -s https://raw.githubusercontent.com/pote/gpm/v1.3.2/bin/gpm | bash  && \
   go get github.com/julienschmidt/httprouter && \
   go get github.com/mattn/go-sqlite3 && \
   go build 

EXPOSE 35800

#run ls /go/bin/ 
RUN ls /go/src/github.com/asiainfoLDP/datahub/datahub
COPY start-notebook.sh /usr/local/bin/
CMD ["start-notebook.sh"]

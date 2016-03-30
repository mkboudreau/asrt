FROM alpine

ENV godir github.com/mkboudreau/asrt
COPY . /go/src/$godir

RUN apk update \
	&& apk add make go git curl ca-certificates   \
   	&& export GOPATH=/go  \	
	&& cd /go/src/$godir \
	&& make docker-static-linux \
	&& CGO_ENABLED=0 go install -a -ldflags '-s' github.com/mkboudreau/asrt  \
	&& mv /go/bin/asrt /usr/local/bin/asrt \
  	&& rm -rf /go*  \
  	&& apk del --purge make go git 

RUN mkdir -p /app
WORKDIR /app
ONBUILD COPY . /app

ENTRYPOINT [ "/usr/local/bin/asrt" ]

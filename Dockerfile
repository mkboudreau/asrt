FROM alpine

RUN mkdir -p /app
WORKDIR /app

RUN apk update \
	&& apk add go git ca-certificates   \
   	&& mkdir /go && export GOPATH=/go  \	
	&& CGO_ENABLED=0 go get -a -ldflags '-s' github.com/mkboudreau/asrt  \
	&& mv /go/bin/asrt /usr/local/bin/asrt \
  	&& rm -rf /go*  \
  	&& apk del --purge go git 

ONBUILD COPY . /app

ENTRYPOINT [ "/usr/local/bin/asrt" ]

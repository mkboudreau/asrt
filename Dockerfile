FROM golang:onbuild

RUN ln -s /go/bin/app /go/bin/asrt

ENTRYPOINT [ "asrt" ]


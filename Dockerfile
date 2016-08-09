FROM golang:1.6.3
MAINTAINER itsyou.online

RUN git clone https://github.com/jteeuwen/go-bindata.git $GOPATH/src/github.com/jteeuwen/go-bindata
WORKDIR $GOPATH/src/github.com/jteeuwen/go-bindata
RUN git checkout a0ff2567cfb70903282db057e799fd826784d41d
RUN go get github.com/jteeuwen/go-bindata/...

RUN git clone https://github.com/Jumpscale/go-raml.git $GOPATH/src/github.com/Jumpscale/go-raml
WORKDIR $GOPATH/src/github.com/Jumpscale/go-raml
RUN git checkout f51adf00a0868be154a67396ad50f6d4d8be7cf0
RUN ./build.sh

ENV CGO_ENABLED 0
WORKDIR /go/src/github.com/itsyouonline/identityserver

EXPOSE 8080

ENTRYPOINT go build && ./identityserver -d

FROM golang:1.6
MAINTAINER itsyou.online

ENV CGO_ENABLED 0
WORKDIR /go/src/github.com/itsyouonline/identityserver

EXPOSE 8080

ENTRYPOINT go build && ./identityserver -d

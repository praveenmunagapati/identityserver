FROM golang
MAINTAINER itsyou.online

ADD . /go/src/github.com/itsyouonline/identityserver

RUN cd /go/src/github.com/itsyouonline/identityserver && go get

RUN go install github.com/itsyouonline/identityserver

EXPOSE 8080

ENTRYPOINT ["/go/bin/identityserver"]

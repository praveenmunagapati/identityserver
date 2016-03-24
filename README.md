# identityserver

[![Build Status](https://travis-ci.org/itsyouonline/identityserver.svg?branch=master)](https://travis-ci.org/itsyouonline/identityserver)

## install and run locally

### Installation

```
go get github.com/itsyouonline/identityserver
```

### Run

```
identityserver
```

To see the available options (like changing the default mongo connectionstring), execute `identityserver -h`.

Browse to https://dev.itsyou.online:8443

* dev.itsyou.online is a public DNS entry that points to 127.0.0.1 and ::1


### Docker-compose

You can run via [Docker-compose](https://docs.docker.com/compose/install/). You will get a running `identityserver` with its own [Mongo](https://hub.docker.com/_/mongo/) docker instance.

```
docker-compose up
```

then browse to https://dev.itsyou.online:8443 (assuming docker runs on localhost)

## Contribute

When you want to contribute to the development, follow the [contribution guidelines](contributing.md).

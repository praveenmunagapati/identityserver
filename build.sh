#!/usr/bin/env bash
set -e

docker build -t itsyouonlinebuilder .
docker run --rm -v "$PWD":/go/src/github.com/itsyouonline/identityserver --entrypoint go  itsyouonlinebuilder build -ldflags '-s' -v -o dist/identityserver
docker build -t itsyouonline/identityserver:0.9-beta-2 -f DockerfileMinimal .

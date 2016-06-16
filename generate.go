package main

//This file contains the go:generate commands

//go-raml https://github.com/Jumpscale/go-raml server code generation from the RAML specification
//TODO: fix serverside code generation with go-raml
//go:generate go-raml client -l go --dir clients/go/itsyouonline --ramlfile specifications/api/itsyouonline.raml
//go:generate go-raml client -l python --dir clients/python/itsyouonline --ramlfile specifications/api/itsyouonline.raml

//go:generate go-bindata -pkg specifications -prefix specifications/api -o specifications/packaged.go specifications/api/...

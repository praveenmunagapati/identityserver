package main

//This file contains the go:generate commands

//go-raml https://github.com/Jumpscale/go-raml server code generation from the RAML specification
//go:generate go-raml server --ramlfile specifications/api/users.raml --dir identityservice/user --package user --no-main
//go:generate go-raml server --ramlfile specifications/api/companies.raml --dir identityservice/company --package company --no-main
//go:generate go-raml server --ramlfile specifications/api/organizations.raml --dir identityservice/organization --package organization --no-main
//go:generate go-raml server --ramlfile specifications/api/contracts.raml --dir identityservice/contracts --package contract --no-main

package http

import (
	"net/http"
)

type Route struct {
	Name        string
	Methods     RouteMethods
	Path        string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

type RouteMethods []string

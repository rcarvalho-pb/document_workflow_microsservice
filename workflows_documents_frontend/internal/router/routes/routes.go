// Package routes defines the Route type and application routes
package routes

import (
	"fmt"
	"net/http"
)

type Route struct {
	URL            string
	Method         string
	Function       func(http.ResponseWriter, *http.Request)
	Authentication bool
}

func ConfigRouter(r *http.ServeMux) *http.ServeMux {
	routes := make([]Route, 0)
	routes = append(routes, GetPagesRoutes()...)

	for _, route := range routes {
		if route.Authentication {
		} else {
			r.HandleFunc(fmt.Sprintf("%s %s", route.Method, route.URL), route.Function)
		}
	}

	return r
}

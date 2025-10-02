// Package router is responsable to create router and give it back
package router

import "net/http"

func NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	return r
}

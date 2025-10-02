package routes

import "net/http"

func GetPagesRoutes() []Route {
	return []Route{
		{
			URL:            "/",
			Method:         http.MethodGet,
			Function:       func(w http.ResponseWriter, r *http.Request) {},
			Authentication: false,
		},
		{
			URL:            "/list-docs",
			Method:         http.MethodGet,
			Function:       func(w http.ResponseWriter, r *http.Request) {},
			Authentication: false,
		},
	}
}

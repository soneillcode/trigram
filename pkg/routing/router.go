package routing

import (
	"net/http"
)

type Router struct {
	routes          []*Route
	notFoundHandler http.Handler
}

func NewRouter() *Router {
	return &Router{
		notFoundHandler: http.HandlerFunc(http.NotFound),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		handler := route.GetHandler(req)
		if handler != nil {
			(*handler)(w, req)
			return
		}
	}
	r.notFoundHandler.ServeHTTP(w, req)
}

func (r *Router) AddRoute(path string, method string, handler Handler) {
	r.routes = append(r.routes, NewRoute(path, method, handler))
}

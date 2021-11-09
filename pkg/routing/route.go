package routing

import (
	"net/http"
)

type Route struct {
	path    string
	method  string
	handler Handler
}

type Match struct {
	Handler Handler
}

func NewRoute(path string, method string, handler Handler) *Route {
	return &Route{
		path:    path,
		method:  method,
		handler: handler,
	}
}

func (r *Route) GetHandler(req *http.Request) *Handler {
	if req.Method != r.method {
		return nil
	}
	if r.path != req.URL.Path {
		return nil
	}
	return &r.handler
}

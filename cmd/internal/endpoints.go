package internal

import (
	"net"
)

type Handler interface {
	ServeHTTP(*Request, *Response)
}

type HandlerFunc func(*Request, *Response)

func (f HandlerFunc) ServeHTTP(req *Request, res *Response) {
	f(req, res)
}

type Router struct {
	conn   net.Conn
	routes map[string]Handler
}

func CreateRouter(c net.Conn) *Router {
	return &Router{
		conn:   c,
		routes: make(map[string]Handler),
	}
}

func (r *Router) AddRoute(path string, handler Handler) {
	r.routes[path] = handler
}

func (r *Router) Serve(req *Request, res *Response) {
	handler, ok := r.routes[req.Params.Path]

	type Result struct {
		Status  string
		Message string
	}

	if !ok {
		result := &Result{
			Status:  "Error",
			Message: "Rota não encontrada",
		}

		res.SendResponse(404, result)
		return
	}

	handler.ServeHTTP(req, res)
}

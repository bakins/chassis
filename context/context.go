// Package Context provides simple http wrappers.
package context

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

type (
	// Handler responders to an HTTP request.
	Handler interface {
		ServeHTTP(*Context)
	}

	// The HandlerFunc type is an adapter.
	HandlerFunc func(*Context)

	// Context represents a single HTTP request/response.
	Context struct {
		Context context.Context
		Writer  ResponseWriter
		Request *Request
		Params  httprouter.Params
	}
)

// HTTPHandler is a middleware that wraps a net/http Handler
func HTTPHandler(h http.Handler) Handler {
	return HandlerFunc(func(ctx *Context) {
		ctx := Context{
			Context: context.Background(),
		}
	})
}

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(ctx *Context) {
	f(ctx)
}

// NewHandler is used to create a context and kick of the serving chain.
func NewHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) Handler {
	return HandlerFunc(func(ctx *Context) {
		ctx := Context{
			Context: context.Background(),
			Writer:  w,
			Request: r,
			Params:  ps,
		}
	})
}

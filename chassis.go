// Package chassis provides simple http wrappers.
package chassis

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	// Router is a simple wrapper around http router
	Router struct {
		Router httprouter.Router
	}

	// ChainConstructor is used for middleware
	ChainConstructor func(Handler) Handler

	// Chain is a representation of middleware
	Chain struct {
		constructors []ChainConstructor
	}
)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(ctx *Context) {
	f(ctx)
}

func NewRouter() *Router {
	r := Router{
		Router: httprouter.New(),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	return r.Router.ServerHTTP(w, r)
}

func (r *Router) Handle(method, path string, h Handler) {
	f := func(w http.ResponseWriter, r *http.Request, ps Params) {
		ctx := Context{
			Context: context.Background(),
			Writer:  w,
			Request: r,
			Params:  ps,
		}
		h(ctx)
	}

	r.Router.Handle(method, path, f)
}

func NewChain(constructors ...ChainConstructor) Chain {
	return Chain{append(([]ChainConstructor)(nil), constructors...)}
}

func (c Chain) Then(h Handler) Handler {
	for i := range c.constructors {
		h = c.constructors[len(c.constructors)-1-i](h)
	}

	return h
}

func (c Chain) Append(constructors ...ChainConstructor) Chain {
	newCons := make([]ChainConstructor, len(c.constructors)+len(constructors))
	copy(newCons, c.constructors)
	copy(newCons[len(c.constructors):], constructors)

	return New(newCons...)
}

func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.constructors...)
}

func (ctx *Context) JSON(code int, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	ctx.Writer.Header().Set("Content-Length", strconv.Itoa(len(data)))
	ctx.Writer.Header().Set("Content-Type", "application/json")
	_, err = ctx.Writer.Write(data)
	return err
}

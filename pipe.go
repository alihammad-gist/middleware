// Package middleware provides abstraction for using middlewares, middleware can be
// thought of as something that sits in the middle of a process. It is a function, that
// when given its own continuation (function) will produce continuation for middlewares
// that depend on it.
package middleware

import (
	"net/http"
)

// NewPipe creates a new middleware pipe. It can be thought of an expression
func NewPipe(mws ...Middleware) *Pipe {

	// source https://github.com/go-midway/midway/blob/master/middleware.go
	middleware := func(inner http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			inner = mws[i].Process(inner)
		}
		return inner
	}

	return &Pipe{
		accum: Func(middleware),
	}
}

// Pipe is a pipeline of middlewares
type Pipe struct {
	accum Middleware // accumulated middleware
}

// ServeHTTP implements http.Handler interface, call this function to
// start the Middleware pipe evaluation process.
func (p *Pipe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// we are piping an noOpHandler to the last middleware
	// in order to produce the actual handler
	handler := p.accum.Process(http.HandlerFunc(NoOpHandler))
	handler.ServeHTTP(w, r)
}

// Process implements Middleware interface so Pipe itself can be
// used as a middleware.
func (p *Pipe) Process(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := p.accum.Process(next)
		handler.ServeHTTP(w, r)
	})
}

// NoOpHandler is used by Pipe for generating an http.Handler
func NoOpHandler(_ http.ResponseWriter, _ *http.Request) {}

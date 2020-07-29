package middleware

import (
	"net/http"
)

// NewPipe creates a new middleware pipe. It can be thought of an expression
func NewPipe(mws ...Middleware) *Pipe {
	flat := make([]Middleware, 0)

	for i := 0; i < len(mws); i++ {
		switch m := mws[i].(type) {
		case *Pipe:
			flat = append(flat, m.middlewares...)
		default:
			flat = append(flat, m)
		}
	}

	return &Pipe{
		middlewares: flat,
	}
}

// Pipe is a pipeline of middlewares
type Pipe struct {
	middlewares []Middleware
	cached      http.Handler
}

func (p *Pipe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.cached == nil {
		p.cached = p.Process(http.HandlerFunc(NoOpHandler))
	}

	p.cached.ServeHTTP(w, r)
}

func (p *Pipe) Process(next http.Handler) http.Handler {
	for i := len(p.middlewares) - 1; i >= 0; i-- {
		next = p.middlewares[i].Process(next)
	}
	return next
}

func NoOpHandler(_ http.ResponseWriter, _ *http.Request) {}

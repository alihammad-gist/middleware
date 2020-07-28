// Package middleware provides abstraction on top of the famous middleware
// pattern: `middleware(http.Handler) http.Handler`. It provide middleware
// pipe and some predefined middlewares.
package middleware

import (
	"net/http"
)

// Middleware is part of the puzzle that is request evaluation. Conceptually
// if request evaluation was represented with an expression, then a middleware
// would be its sub expression.
type Middleware interface {
	Process(http.Handler) http.Handler
}

// Func is a function type that implements Middleware interface. It can be used
// to create middlewares. middleware.Func(func(h http.Handler) http.Handler {...})
type Func func(http.Handler) http.Handler

// Process does not initiate the req/resp evaluation process, rather
// it can be thought of as an expression builder
func (mf Func) Process(h http.Handler) http.Handler {
	return mf(h)
}

// Terminate will convert a handler to a terminating middleware that
// will ignore rest of the continuation. Its technically called Identity continuation,
// where the continuation of an expression is the expression itself. For an exp: 6
// continuation: (lambda (k) (k))
func Terminate(h http.Handler) Middleware {
	return Func(
		// does not take any continuation
		func(_ http.Handler) http.Handler {
			return h
		},
	)
}

// Continue will execute rest of the continuation after the http.Handler
// finishes its execution
func Continue(h http.Handler) Middleware {
	return Func(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.ServeHTTP(w, r)
				next.ServeHTTP(w, r)
			})
		},
	)
}

package middleware_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/alihammad-gist/middleware"
)

const (
	middlewareLabel = "middleware => %s\n"
	handlerLabel    = "handler => %s\n"
)

func createMiddleware(label string, t *testing.T) middleware.Middleware {
	return middleware.Func(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(createMiddlewareLabel(label)))
			t.Log(createMiddlewareLabel(label))
			next.ServeHTTP(w, r)
		})
	})
}

func createHandler(label string, t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(createHandlerLabel(label)))
		t.Log(createHandlerLabel(label))
	})
}

func createMiddlewareLabel(label string) string {
	return fmt.Sprintf(middlewareLabel, label)
}

func createHandlerLabel(label string) string {
	return fmt.Sprintf(handlerLabel, label)
}

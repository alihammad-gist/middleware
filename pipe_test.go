package middleware_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alihammad-gist/middleware"
)

func TestSomething(t *testing.T) {
	p := middleware.NewPipe(
		middleware.NewPipe(
			createMiddleware("First", t),
			createMiddleware("Second", t),
		),
		createMiddleware("Third", t),
		middleware.NewPipe(
			createMiddleware("Fourth", t),
			createMiddleware("Fifth", t),
			createMiddleware("Sixth", t),
			middleware.NewPipe(
				createMiddleware("Seventh", t),
				createMiddleware("Eigth", t),
				createMiddleware("Nineth", t),
			),
		),
		middleware.Terminate(createHandler("Tenth", t)),
		middleware.NewPipe(
			createMiddleware("Eleventh", t),
			createMiddleware("...", t),
		),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://www.example.com/path", nil)

	p.ServeHTTP(w, r)

	body, err := ioutil.ReadAll(w.Result().Body)

	if err != nil {
		// error reading content of response body
		t.Fatal(err)
	}

	if string(body) != strings.Join(
		[]string{
			createMiddlewareLabel("First"),
			createMiddlewareLabel("Second"),
			createMiddlewareLabel("Third"),
			createMiddlewareLabel("Fourth"),
			createMiddlewareLabel("Fifth"),
			createMiddlewareLabel("Sixth"),
			createMiddlewareLabel("Seventh"),
			createMiddlewareLabel("Eigth"),
			createMiddlewareLabel("Nineth"),
			createHandlerLabel("Tenth"),
		},
		"",
	) {
		t.Fatal("Invalid http response")
		t.Fail()
	}

}

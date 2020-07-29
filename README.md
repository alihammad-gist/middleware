# Middleware

This package provides middleware pipe abstraction on top of the famous `func(http.Handler) http.Handler` middleware type, it works well with routing packages like https://github.com/bmizerany/pat . General `Continue` and `Terminate` wrappers for `http.Handler` are also included to reduce boilerplate code. `Continue` covers the frequent case of executing the next in line `http.Handler` at the end of  the current one (you can just pass the current one to `Continue`), `Terminate` will just ignore executing the rest of the `http.Handler`s. The package also includes following predefined middlewares.

1. `RequestID` assigns a unique (int64) ID to each request (context).
2. `RequestLogger` logs each http request.

Any middlewares with signature `func(http.Handler) http.Handler` can be used with this package. Following are a few examples of programs using this package.

```golang
package main

import (
    "log"
    "net/http"
    "os"

    "github.com/alihammad-gist/middleware"
    "github.com/bmizerany/pat"
)

func main() {
    mux := pat.New()

    loggerOpts := middleware.LoggerOptions{
        RequestID:  true,
        RemoteAddr: true,
        RequestURI: true,
        Host:       true,
    }

    http.Handle("/", middleware.NewPipe(
        middleware.RequestID(),
        middleware.RequestLogger(log.New(os.Stdout, "HTTP-Request ", log.LstdFlags), loggerOpts),
        middleware.Terminate(mux), // attach your pat multiplexer
    ))

    log.Println("About to listen on localhost:9999")
    if err := http.ListenAndServe(":9999", nil); err != nil {
        log.Fatal(err)
    }
}
```
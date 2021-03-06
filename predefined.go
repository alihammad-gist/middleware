package middleware

import (
	"context"
	"github.com/edwingeng/wuid/callback/wuid"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	// KeyRequestID is the key in request context to
	KeyRequestID struct{}

	// LoggerOptions is used by RequestLoggerMiddleware factory
	LoggerOptions struct {
		RequestID  bool
		RemoteAddr bool
		RequestURI bool
		Host       bool
	}
)

var (
	// universal ID generator used by requestIDMiddleware
	uidGen *wuid.WUID
)

func init() {
	uidGen = wuid.NewWUID("KeyRequestID", nil)
	uidGen.LoadH28WithCallback(func() (int64, func(), error) {
		return time.Now().UnixNano(), nil, nil
	})
}

// RequestID assigns a unique ID (int64) to each request (context). The
// ID can accessed from the context using middleware.KeyRequestID
func RequestID() Middleware {
	return Func(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, KeyRequestID{}, uidGen.Next())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

// RequestLogger logs each http request, logging the specified fields in passed in
// middleware.LoggerOptions value.
func RequestLogger(logger *log.Logger, opts LoggerOptions) Middleware {
	return Continue(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		msgs := []string{}

		if opts.RequestID {
			reqID := r.Context().Value(KeyRequestID{})
			if reqID != nil {
				id := reqID.(int64)
				msgs = append(msgs, "(RequestID "+strconv.FormatInt(id, 10)+")")
			} else {
				msgs = append(msgs, "(RequestID N/A)")
				log.Println("Consider adding middleware.RequestId to your root middleware Pipe")
			}
		}

		if opts.RemoteAddr {
			msgs = append(msgs, "(RemoteAddr "+r.RemoteAddr+")")
		}

		if opts.RequestURI {
			msgs = append(msgs, "(RequestURI "+r.RequestURI+")")
		}

		if opts.Host {
			msgs = append(msgs, "(Host "+r.Host+")")
		}

		if len(msgs) > 0 {
			logger.Println(strings.Join(msgs, " "))
		}
	}))
}

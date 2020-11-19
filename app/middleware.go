package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type middleware func(next http.Handler) http.Handler

// middlewareChain takes an array of middleware functions and a final handler
// and chains them in the same order as in the array.
// middlewareChain([]middleware{m1, m2, m3}, h) ==> m1(m2(m3(h)))
func middlewareChain(funcs []middleware, final http.Handler) http.Handler {
	for i := range funcs {
		final = funcs[len(funcs)-1-i](final)
	}
	return final
}

func trimSuffixMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := statusRecorder{w, http.StatusOK}
		next.ServeHTTP(&rec, r)
		log.WithFields(log.Fields{
			"req":    fmt.Sprintf("%s %s", r.Method, r.RequestURI),
			"status": rec.status,
		}).Info("handled request")
	})
}

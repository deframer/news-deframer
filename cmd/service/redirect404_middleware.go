package main

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/deframer/news-deframer/pkg/config"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
	headers    http.Header
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		headers:        make(http.Header),
	}
}

func (rw *responseRecorder) Header() http.Header {
	return rw.headers
}

func (rw *responseRecorder) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

func (rw *responseRecorder) Write(b []byte) (int, error) {
	return rw.body.Write(b)
}

func isJSONRequest(r *http.Request) bool {
	accept := strings.ToLower(r.Header.Get("Accept"))
	return strings.Contains(accept, "application/json") || strings.Contains(accept, "application/*+json")
}

func redirect404Middleware(next http.Handler) http.Handler {
	cfg, err := config.Load()
	if err != nil || strings.TrimSpace(cfg.RedirectWebRequest404URL) == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isJSONRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		rw := newResponseRecorder(w)
		next.ServeHTTP(rw, r)

		if rw.statusCode == http.StatusNotFound {
			http.Redirect(w, r, cfg.RedirectWebRequest404URL, http.StatusFound)
			return
		}

		for key, values := range rw.headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(rw.statusCode)
		if rw.body.Len() == 0 {
			return
		}
		_, _ = w.Write(rw.body.Bytes())
	})
}

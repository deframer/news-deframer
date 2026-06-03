package main

import (
	"net/http"
	"strings"

	"github.com/deframer/news-deframer/pkg/config"
)

func corsMiddleware(next http.Handler) http.Handler {
	cfg, err := config.Load()
	if err != nil || strings.TrimSpace(cfg.CORSAllowedOrigins) == "" {
		// If cors_allowed is empty, do nothing.
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD")
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.Header().Set("Vary", "Origin")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

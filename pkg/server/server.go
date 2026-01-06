package server

import (
	"context"
	"net"
	"net/http"

	"github.com/egandro/news-deframer/pkg/config"
)

type Server struct {
	httpServer *http.Server
}

func New(ctx context.Context, cfg *config.Config) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/rss", func(w http.ResponseWriter, r *http.Request) {
		// The requirement is to "accept and return xml only".
		// We enforce the response Content-Type to XML.
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
<channel>
  <title>News Deframer</title>
</channel>
</rss>`))
	})

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: mux,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

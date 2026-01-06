package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/egandro/news-deframer/pkg/config"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

type RSSRequest struct {
	URL      string
	Lang     string
	MaxScore float64
	Embedded bool
}

func New(ctx context.Context, cfg *config.Config) *Server {
	s := &Server{
		logger: slog.With("component", "server"),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", s.handlePing)
	mux.HandleFunc("/rss", s.handleRSS)

	mux.HandleFunc("/site", s.handleSite)
	mux.HandleFunc("/item", s.handleItem)

	s.httpServer = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	return s
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("ping")
	w.Write([]byte("pong"))
}

func (s *Server) handleRSS(w http.ResponseWriter, r *http.Request) {
	// curl "http://localhost:8080/rss?url=http%3A%2F%2Fexample.com%2Ffeed.xml&lang=en&max_score=0.5&embedded=true"
	q := r.URL.Query()
	req := RSSRequest{
		URL:  q.Get("url"),
		Lang: q.Get("lang"),
	}

	if req.URL == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		s.logger.Debug("invalid url", "error", err)
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	if v, err := strconv.ParseFloat(q.Get("max_score"), 64); err == nil {
		req.MaxScore = v
	}
	if v, err := strconv.ParseBool(q.Get("embedded")); err == nil {
		req.Embedded = v
	}

	// The requirement is to "accept and return xml only".
	// We enforce the response Content-Type to XML.
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0" xmlns:debug="http://news-deframer/debug">
<channel>
  <title>News Deframer</title>
  <debug:source_url>%s</debug:source_url>
  <debug:lang>%s</debug:lang>
  <debug:max_score>%f</debug:max_score>
  <debug:embedded>%t</debug:embedded>
</channel>
</rss>`, req.URL, req.Lang, req.MaxScore, req.Embedded)
}

func (s *Server) handleSite(w http.ResponseWriter, r *http.Request) {
	// curl "http://localhost:8080/site?url=http%3A%2F%2Fexample.com"
	q := r.URL.Query()
	reqURL := q.Get("url")

	if reqURL == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(reqURL); err != nil {
		s.logger.Debug("invalid url", "error", err)
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": reqURL})
}

func (s *Server) handleItem(w http.ResponseWriter, r *http.Request) {
	// curl "http://localhost:8080/item?url=http%3A%2F%2Fexample.com%2Farticle"
	q := r.URL.Query()
	reqURL := q.Get("url")

	if reqURL == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(reqURL); err != nil {
		s.logger.Debug("invalid url", "error", err)
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": reqURL})
}

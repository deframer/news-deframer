package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/facade"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	facade     facade.Facade
}

type RSSRequest struct {
	URL      string
	Lang     string
	MaxScore float64
	Embedded bool
}

func New(ctx context.Context, cfg *config.Config, f facade.Facade) *Server {
	s := &Server{
		logger: slog.With("component", "server"),
		facade: f,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", s.handlePing)
	mux.HandleFunc("GET /hostname", s.handleHostname)
	mux.HandleFunc("/rss", s.handleRSSProxy)

	mux.HandleFunc("/site", s.handleSite)
	mux.HandleFunc("/item", s.handleItem)

	s.httpServer = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ReadHeaderTimeout: 5 * time.Second,
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
	if _, err := w.Write([]byte("pong")); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

func (s *Server) handleHostname(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		s.logger.Error("failed to get hostname", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"hostname": hostname}); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

func (s *Server) handleRSSProxy(w http.ResponseWriter, r *http.Request) {
	// curl "http://localhost:8080/rss?url=http%3A%2F%2Fdummy&lang=en&max_score=0.5&embedded=true"
	// curl "http://localhost:8080/rss?url=http%3A%2F%2Fwordpress%2Ffeed"
	// curl "http://localhost:8080/rss?url=http%3A%2F%2Flocalhost%3A8003%2Ffeed"
	// inside freshrss http://service:8080/rss?url=http%3A%2F%2Fwordpress%2Ffeed
	q := r.URL.Query()
	req := RSSRequest{
		URL:  strings.TrimSuffix(q.Get("url"), "/"),
		Lang: q.Get("lang"),
	}

	if req.URL == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}

	u, err := url.ParseRequestURI(req.URL)
	if err != nil {
		s.logger.Debug("invalid url", "error", err)
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	hasFeed, err := s.facade.HasFeedByUrl(r.Context(), u)
	if err != nil || !hasFeed {
		if err != nil {
			s.logger.Error("HasFeed failed", "error", err)
		} else {
			s.logger.Debug("feed not found", "url", req.URL)
		}
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if v, err := strconv.ParseFloat(q.Get("max_score"), 64); err == nil {
		req.MaxScore = v
	}
	if v, err := strconv.ParseBool(q.Get("embedded")); err == nil {
		req.Embedded = v
	}

	filter := facade.RSSProxyFilter{
		URL:      req.URL,
		Lang:     req.Lang,
		MaxScore: req.MaxScore,
		Embedded: req.Embedded,
	}

	xmlData, err := s.facade.GetRssProxyFeed(r.Context(), &filter)
	if err != nil {
		s.logger.Error("failed to get rss proxy feed", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	if _, err := w.Write([]byte(xmlData)); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

func (s *Server) handleSite(w http.ResponseWriter, r *http.Request) {
	// curl "http://localhost:8080/site?url=http%3A%2F%2Fexample.com"
	q := r.URL.Query()
	reqURL := strings.TrimSuffix(q.Get("url"), "/")

	if reqURL == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}

	u, err := url.ParseRequestURI(reqURL)
	if err != nil {
		s.logger.Debug("invalid url", "error", err)
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	hasFeed, err := s.facade.HasFeedByUrl(r.Context(), u)
	if err != nil || !hasFeed {
		if err != nil {
			s.logger.Error("HasFeed failed", "error", err)
		} else {
			s.logger.Debug("feed not found", "url", reqURL)
		}
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"url":      reqURL,
		"has_feed": hasFeed,
	}); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

func (s *Server) handleItem(w http.ResponseWriter, r *http.Request) {
	// 0 entries
	// curl "http://localhost:8080/item?url=http%3A%2F%2Fexample.com%2Farticle"
	// 2 entries
	// curl "http://localhost:8080/item?url=http%3A%2F%2Fdummy-enforced%2Fitem-1"
	// 1 entry
	// curl "http://localhost:8080/item?url=http%3A%2F%2Fdummy-open%2Fitem-2"

	q := r.URL.Query()
	reqURL := strings.TrimSuffix(q.Get("url"), "/")

	if reqURL == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}

	u, err := url.ParseRequestURI(reqURL)
	if err != nil {
		s.logger.Debug("invalid url", "error", err)
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	items, err := s.facade.GetItems(r.Context(), u)
	if err != nil || len(items) == 0 {
		if err != nil {
			s.logger.Error("GetItems failed", "error", err)
		} else {
			s.logger.Debug("no items found", "url", reqURL)
		}
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

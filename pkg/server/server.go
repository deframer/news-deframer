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
	Max      float64
	Embedded bool
}

func New(ctx context.Context, cfg *config.Config, f facade.Facade) *Server {
	s := &Server{
		logger: slog.With("component", "server"),
		facade: f,
	}

	mux := http.NewServeMux()

	// debug stuff
	mux.HandleFunc("/ping", s.handlePing)
	mux.HandleFunc("/hostname", s.handleHostname)

	// Middleware
	protect := func(h http.HandlerFunc) http.Handler {
		if cfg.BasicAuthUser != "" && cfg.BasicAuthPassword != "" {
			return s.basicAuthMiddleware(cfg.BasicAuthUser, cfg.BasicAuthPassword, h)
		}
		return h
	}

	if cfg.BasicAuthUser != "" && cfg.BasicAuthPassword != "" {
		s.logger.Info("Basic auth enabled")
	}

	mux.Handle("/rss", protect(s.handleRSSProxy))
	mux.Handle("/api/item", protect(s.handleItem))
	mux.Handle("/api/site", protect(s.handleSite))
	mux.Handle("/api/domains", protect(s.handleDomains))

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

func (s *Server) basicAuthMiddleware(user, password string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != user || p != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
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
	q := r.URL.Query()
	req := RSSRequest{
		URL:  strings.TrimSuffix(q.Get("url"), "/"),
		Lang: q.Get("lang"),
	}

	if req.URL == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}

	if v, err := strconv.ParseFloat(q.Get("max_score"), 64); err == nil {
		req.Max = v
	}
	if v, err := strconv.ParseBool(q.Get("embedded")); err == nil {
		req.Embedded = v
	}

	filter := facade.RSSProxyFilter{
		URL:      req.URL,
		Lang:     req.Lang,
		Max:      req.Max,
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

func (s *Server) handleItem(w http.ResponseWriter, r *http.Request) {
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

	item, err := s.facade.GetFirstItemForUrl(r.Context(), u)
	if err != nil {
		s.logger.Error("failed to get item", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if item == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

func (s *Server) handleSite(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	rootDomain := strings.TrimSuffix(q.Get("root"), "/")

	if rootDomain == "" {
		http.Error(w, "missing root", http.StatusBadRequest)
		return
	}

	var maxScore float64
	if v, err := strconv.ParseFloat(q.Get("max_score"), 64); err == nil {
		maxScore = v
	}

	items, err := s.facade.GetItemsForRootDomain(r.Context(), rootDomain, maxScore)
	if err != nil || len(items) == 0 {
		if err != nil {
			s.logger.Error("GetItemsForRootDomain failed", "error", err)
		} else {
			s.logger.Debug("no items found", "url", rootDomain)
		}
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

func (s *Server) handleDomains(w http.ResponseWriter, r *http.Request) {
	domains, err := s.facade.GetRootDomains(r.Context())
	if err != nil {
		s.logger.Error("failed to get root domains", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(domains); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

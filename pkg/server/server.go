package server

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/facade"
)

// responseWriter is a wrapper for http.ResponseWriter to capture the status code and body
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.body.Write(b)
}

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	facade     facade.Facade
	cfg        *config.Config
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
		cfg:    cfg,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", s.handlePing)
	mux.HandleFunc("/hostname", s.handleHostname)
	mux.Handle("/rss", s.basicAuthMiddleware(s.handleRSSProxy))
	mux.Handle("/api/item", s.basicAuthMiddleware(s.handleItem))
	mux.Handle("/api/site", s.basicAuthMiddleware(s.handleSite))
	mux.Handle("/api/domains", s.basicAuthMiddleware(s.handleDomains))
	mux.Handle("/api/trends/topbydomain", s.basicAuthMiddleware(s.handleTopTrendsByDomain))
	mux.Handle("/api/trends/contextbydomain", s.basicAuthMiddleware(s.handleContextByDomain))
	mux.Handle("/api/trends/lifecyclebydomain", s.basicAuthMiddleware(s.handleLifecycleByDomain))

	// Chain middlewares: etag -> cache-control -> mux
	var handler http.Handler = mux
	if !cfg.DisableETag {
		handler = s.etagMiddleware(s.cacheControlMiddleware(handler))
	}

	s.httpServer = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ReadHeaderTimeout: 5 * time.Second,
	}
	return s
}

func (s *Server) cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			maxAge := int(config.PollingInterval.Seconds() / 2)
			if strings.HasPrefix(r.URL.Path, "/rss") {
				w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
			} else if strings.HasPrefix(r.URL.Path, "/api/") {
				w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(config.ETagTTL.Seconds())))
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) etagMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		// After handler has run, we can inspect the body and status code

		// Don't set ETag for empty bodies, just write headers and return
		if rw.body.Len() == 0 {
			w.WriteHeader(rw.statusCode)
			return
		}

		hash := sha256.Sum256(rw.body.Bytes())
		etag := `"` + hex.EncodeToString(hash[:]) + `"`
		w.Header().Set("ETag", etag)

		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, etag) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		w.WriteHeader(rw.statusCode)
		if _, err := w.Write(rw.body.Bytes()); err != nil {
			s.logger.Error("failed to write response", "error", err)
		}
	})
}

func (s *Server) basicAuthMiddleware(next http.HandlerFunc) http.Handler {
	if s.cfg.BasicAuthUser == "" || s.cfg.BasicAuthPassword == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != s.cfg.BasicAuthUser || p != s.cfg.BasicAuthPassword {
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

func (s *Server) handleTopTrendsByDomain(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("handleTopTrendsByDomain", "url", r.URL.String())
	q := r.URL.Query()

	domain := q.Get("domain")
	if domain == "" {
		domain = q.Get("d")
	}

	lang := q.Get("lang")
	if lang == "" {
		lang = q.Get("l")
	}

	if domain == "" || lang == "" {
		http.Error(w, "missing domain or lang", http.StatusBadRequest)
		return
	}

	days := 7
	if v, err := strconv.Atoi(q.Get("days")); err == nil {
		days = v
	}

	trends, err := s.facade.GetTopTrendByDomain(r.Context(), domain, lang, days)
	if err != nil {
		s.logger.Error("failed to get top trends", "error", err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(trends); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

func (s *Server) handleContextByDomain(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("handleContextByDomain", "url", r.URL.String())
	q := r.URL.Query()

	term := q.Get("term")
	if term == "" {
		term = q.Get("t")
	}

	domain := q.Get("domain")
	if domain == "" {
		domain = q.Get("d")
	}

	lang := q.Get("lang")
	if lang == "" {
		lang = q.Get("l")
	}

	if term == "" || domain == "" || lang == "" {
		http.Error(w, "missing term, domain or lang", http.StatusBadRequest)
		return
	}

	days := 7
	if v, err := strconv.Atoi(q.Get("days")); err == nil {
		days = v
	}

	contexts, err := s.facade.GetContextByDomain(r.Context(), term, domain, lang, days)
	if err != nil {
		s.logger.Error("failed to get context by domain", "error", err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(contexts); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

func (s *Server) handleLifecycleByDomain(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("handleLifecycleByDomain", "url", r.URL.String())
	q := r.URL.Query()

	term := q.Get("term")
	if term == "" {
		term = q.Get("t")
	}

	domain := q.Get("domain")
	if domain == "" {
		domain = q.Get("d")
	}

	lang := q.Get("lang")
	if lang == "" {
		lang = q.Get("l")
	}

	if term == "" || domain == "" || lang == "" {
		http.Error(w, "missing term, domain or lang", http.StatusBadRequest)
		return
	}

	days := 7
	if v, err := strconv.Atoi(q.Get("days")); err == nil {
		days = v
	}

	lifecycles, err := s.facade.GetLifecycleByDomain(r.Context(), term, domain, lang, days)
	if err != nil {
		s.logger.Error("failed to get lifecycle by domain", "error", err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lifecycles); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

func (s *Server) handleHostname(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		s.logger.Error("failed to get hostname", "error", err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"hostname": hostname}); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

func (s *Server) handleRSSProxy(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("handleRSSProxy", "url", r.URL.String())
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
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	if _, err := w.Write([]byte(xmlData)); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

func (s *Server) handleItem(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("handleItem", "url", r.URL.String())
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
		http.Error(w, "not found", http.StatusNotFound)
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
	s.logger.Debug("handleSite", "url", r.URL.String())
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
	s.logger.Debug("handleDomains")
	domains, err := s.facade.GetRootDomains(r.Context())
	if err != nil {
		s.logger.Error("failed to get root domains", "error", err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(domains); err != nil {
		s.logger.Error("failed to write response", "error", err)
	}
}

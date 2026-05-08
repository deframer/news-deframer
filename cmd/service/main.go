package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"html"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/facade"
	applog "github.com/deframer/news-deframer/pkg/logger"
	"github.com/deframer/news-deframer/pkg/server"
	"goa.design/clue/log"
)

func main() {
	disableETagFlag := flag.Bool("disable-etag", false, "Disable ETag caching")
	flag.Usage = func() {
		// #nosec G705: usage string is escaped before printing
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", html.EscapeString(filepath.Base(os.Args[0])))
		flag.PrintDefaults()
	}
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf(context.Background(), err, "Failed to load config")
		os.Exit(1)
	}

	disableETag := *disableETagFlag
	logCtx := applog.NewLoggerContext(context.Background(), false)

	hostname, _ := os.Hostname()
	log.Print(logCtx, log.KV{K: "component", V: "service"}, log.KV{K: "hostname", V: hostname}, log.KV{K: "etag_disabled", V: disableETag})
	if cfg.CORSAllowedOrigins != "" {
		log.Print(logCtx, log.KV{K: "message", V: "CORS enabled"}, log.KV{K: "origins", V: cfg.CORSAllowedOrigins})
	}

	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo, err := database.NewRepository(logCtx, cfg)
	if err != nil {
		log.Fatalf(logCtx, err, "Failed to connect to database")
		os.Exit(1)
	}

	f := facade.New(appCtx, cfg, repo)

	srv := server.New(appCtx, cfg, f)

	go func() {
		log.Print(logCtx, log.KV{K: "message", V: "Starting server"}, log.KV{K: "port", V: cfg.Port})
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf(logCtx, err, "Server failed")
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Print(logCtx, log.KV{K: "message", V: "Shutting down..."})

	// Shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		log.Fatalf(logCtx, err, "Server shutdown failed")
	}
	log.Print(logCtx, log.KV{K: "message", V: "Server stopped"})
}

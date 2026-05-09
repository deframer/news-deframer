package main

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"

	infrasvr "github.com/deframer/news-deframer/gen/http/infra/server"
	mobilesvr "github.com/deframer/news-deframer/gen/http/mobile/server"
	openapisvr "github.com/deframer/news-deframer/gen/http/openapi/server"
	rsssvr "github.com/deframer/news-deframer/gen/http/rss/server"
	websvr "github.com/deframer/news-deframer/gen/http/web/server"
	infra "github.com/deframer/news-deframer/gen/infra"
	mobile "github.com/deframer/news-deframer/gen/mobile"
	openapi "github.com/deframer/news-deframer/gen/openapi"
	rss "github.com/deframer/news-deframer/gen/rss"
	web "github.com/deframer/news-deframer/gen/web"
	"goa.design/clue/debug"
	"goa.design/clue/log"
	goahttp "goa.design/goa/v3/http"
)

// handleHTTPServer starts configures and starts a HTTP server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleHTTPServer(ctx context.Context, u *url.URL, openapiEndpoints *openapi.Endpoints, infraEndpoints *infra.Endpoints, rssEndpoints *rss.Endpoints, mobileEndpoints *mobile.Endpoints, webEndpoints *web.Endpoints, wg *sync.WaitGroup, errc chan error, dbg bool) {

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and mount debug and profiler
	// endpoints in debug mode.
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
		if dbg {
			// Mount pprof handlers for memory profiling under /debug/pprof.
			debug.MountPprofHandlers(debug.Adapt(mux))
			// Mount /debug endpoint to enable or disable debug logs at runtime.
			debug.MountDebugLogEnabler(debug.Adapt(mux))
		}
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	var (
		openapiServer *openapisvr.Server
		infraServer   *infrasvr.Server
		rssServer     *rsssvr.Server
		mobileServer  *mobilesvr.Server
		webServer     *websvr.Server
	)
	{
		eh := errorHandler(ctx)
		openapiServer = openapisvr.New(openapiEndpoints, mux, dec, enc, eh, nil)
		infraServer = infrasvr.New(infraEndpoints, mux, dec, enc, eh, nil)
		rssServer = rsssvr.New(rssEndpoints, mux, dec, enc, eh, nil)
		mobileServer = mobilesvr.New(mobileEndpoints, mux, dec, enc, eh, nil)
		webServer = websvr.New(webEndpoints, mux, dec, enc, eh, nil)
	}

	// Configure the mux.
	openapisvr.Mount(mux, openapiServer)
	infrasvr.Mount(mux, infraServer)
	rsssvr.Mount(mux, rssServer)
	mobilesvr.Mount(mux, mobileServer)
	websvr.Mount(mux, webServer)

	var handler http.Handler = mux
	if dbg {
		// Log query and response bodies if debug logs are enabled.
		handler = debug.HTTP()(handler)
	}
	// skip pings
	var noLogRegexp = regexp.MustCompile(`^/(healthz|livez|metrics|ping)$`)
	handler = log.HTTP(ctx, log.WithPathFilter(noLogRegexp))(handler)
	// handler = log.HTTP(ctx)(handler)

	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: u.Host, Handler: handler, ReadHeaderTimeout: time.Second * 60}
	for _, m := range openapiServer.Mounts {
		log.Printf(ctx, "HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range infraServer.Mounts {
		log.Printf(ctx, "HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range rssServer.Mounts {
		log.Printf(ctx, "HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range mobileServer.Mounts {
		log.Printf(ctx, "HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range webServer.Mounts {
		log.Printf(ctx, "HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start HTTP server in a separate goroutine.
		go func() {
			log.Printf(ctx, "HTTP server listening on %q", u.Host)
			errc <- srv.ListenAndServe()
		}()

		<-ctx.Done()
		log.Printf(ctx, "shutting down HTTP server at %q", u.Host)

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf(ctx, "failed to shutdown: %v", err)
		}
	}()
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler(logCtx context.Context) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		log.Printf(logCtx, "ERROR: %s", err.Error())
	}
}

package restserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ilazutin/dadataproxy_go/internal/api/rest_server/handlers/cache"
	"github.com/ilazutin/dadataproxy_go/internal/api/rest_server/handlers/clean"
	"github.com/ilazutin/dadataproxy_go/internal/api/rest_server/handlers/iplocate"
	"github.com/ilazutin/dadataproxy_go/internal/api/rest_server/handlers/migrate"
	"github.com/ilazutin/dadataproxy_go/internal/api/rest_server/handlers/suggest"
	mwLogger "github.com/ilazutin/dadataproxy_go/internal/api/rest_server/middleware/logger"
	service "github.com/ilazutin/dadataproxy_go/internal/service"
)

type ProxyServer struct {
	server *http.Server
	logger *slog.Logger
}

func New(address string, timeout time.Duration, idleTimeout time.Duration, proxyService *service.ProxyService, logger *slog.Logger) *ProxyServer {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/suggest/*", suggest.New(proxyService, logger))

	router.Post("/iplocate/*", iplocate.New(proxyService, logger))

	router.Post("/clean/*", clean.New(proxyService, logger))

	router.Post("/findById/*", suggest.New(proxyService, logger))

	router.Post("/geolocate/*", suggest.New(proxyService, logger))

	router.Post("/cache", cache.New(proxyService, logger))

	router.Post("/migrate", migrate.New(proxyService, logger))

	return &ProxyServer{
		server: &http.Server{
			Addr:         address,
			Handler:      router,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
			IdleTimeout:  idleTimeout,
			// gosec issue: G112: Potential Slowloris Attack.
			ReadHeaderTimeout: 10 * time.Second, //nolint:gomnd
		},
		logger: logger,
	}
}

func (o *ProxyServer) Run(context.Context) error {
	o.logger.Info("Proxy server is running", slog.String("address", o.server.Addr))
	return o.server.ListenAndServe()
}

func (o *ProxyServer) Shutdown(err error) {
	o.logger.Info("Proxy server was interrupted", slog.String("error", err.Error()))
	if shutdownErr := o.server.Shutdown(context.Background()); shutdownErr != nil {
		o.logger.Error("Proxy server shutdown error", slog.String("error", err.Error()))
	}
}

package migrate

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilazutin/dadataproxy_go/internal/api/helper"
	"github.com/ilazutin/dadataproxy_go/internal/service"
)

func New(proxyService *service.ProxyService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cache"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("path", r.URL.Path),
		)

		err := proxyService.MigrateCache(log)
		if err != nil {
			log.Error("Internal server error", slog.String("err", err.Error()))
			helper.ResponseInternalError(w, err)

			return
		}

		helper.ResponseOk(w, nil)
	}
}

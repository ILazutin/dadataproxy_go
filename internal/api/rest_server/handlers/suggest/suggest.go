package suggest

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilazutin/dadataproxy_go/internal/api/helper"
	"github.com/ilazutin/dadataproxy_go/internal/service"
)

func New(proxyService *service.ProxyService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.suggest"

		ignoreCache := r.FormValue("ignore_cache") == "true"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("path", r.URL.Path),
			slog.Bool("ignore_cache", ignoreCache),
		)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("Cannot read body", slog.String("err", err.Error()))
			helper.ResponseErrors(w, err)
			return
		}

		log.Info("Decoded request body", slog.String("request", string(body)))

		result, err := proxyService.SuggestValue(r.URL.Path, string(body), ignoreCache, log)
		if err != nil {
			log.Error("Internal server error", slog.String("err", err.Error()))
			helper.ResponseInternalError(w, err)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		helper.ResponseOk(w, result)
	}
}

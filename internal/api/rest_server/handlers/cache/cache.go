package cache

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilazutin/dadataproxy_go/internal/api/helper"
	"github.com/ilazutin/dadataproxy_go/internal/service"
)

type CacheRequest struct {
	Path  string      `json:"path,omitempty"`
	Query interface{} `json:"query"`
	Body  interface{} `json:"body"`
}

func New(proxyService *service.ProxyService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cache"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("path", r.URL.Path),
		)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("Cannot read body", slog.String("err", err.Error()))
			helper.ResponseErrors(w, err)
			return
		}

		log.Info("Decoded request body", slog.String("request", string(body)))

		var requestBody *CacheRequest
		err = json.Unmarshal(body, &requestBody)

		path := "/clean/address"
		if requestBody.Path != "" {
			path = requestBody.Path
		}

		err = proxyService.SaveToCache(path, requestBody.Query, requestBody.Body, log)
		if err != nil {
			log.Error("Internal server error", slog.String("err", err.Error()))
			helper.ResponseInternalError(w, err)

			return
		}

		helper.ResponseOk(w, nil)
	}
}

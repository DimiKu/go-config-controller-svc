package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"go-config-controller-svc/dto/server_dto"
	"go-config-controller-svc/internal/entities"
	"go-config-controller-svc/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

func DeleteConfigHandler(service ServerService, log *zap.Logger, ctx context.Context) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		var config server_dto.HTTPConfigDto

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(buf.Bytes(), &config)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if err := service.DeleteConfig(ctx, entities.ServerConfig(config)); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		utils.MakeConfigResponse(rw, config)
	}
}

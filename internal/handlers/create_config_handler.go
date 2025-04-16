package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"go-config-controller-svc/dto/server_dto"
	"go-config-controller-svc/internal/entities"
	"go.uber.org/zap"
	"net/http"
)

type ServerService interface {
	CreateConfig(ctx context.Context, config entities.ServerConfig) error
}

func CreateConfigHandler(service ServerService, log *zap.Logger, ctx context.Context) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		var config server_dto.HTTPConfigDto
		var transferConfig entities.ServerConfig

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

		transferConfig.ConfigValue = config.ConfigValue
		transferConfig.ConfigName = config.ConfigName
		transferConfig.ConfigBranch = config.ConfigBranch

		if err := service.CreateConfig(ctx, transferConfig); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

	}
}

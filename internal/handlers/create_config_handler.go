//go:generate mockgen -source=./create_config_handler.go -destination=./create_config_handler_mock.go -package=handlers
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go-config-controller-svc/dto/server_dto"
	"go-config-controller-svc/internal/entities"
	"go-config-controller-svc/internal/interfaces"
	"go.uber.org/zap"
	"net/http"
)

func CreateConfigHandler(service interfaces.ServerService, log *zap.Logger, ctx context.Context) func(rw http.ResponseWriter, r *http.Request) {
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

		if err := validateConf(config); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		transferConfig.ConfigValue = config.ConfigValue
		transferConfig.ConfigName = config.ConfigName
		transferConfig.ConfigBranch = config.ConfigBranch

		log.Info("Create config", zap.String("Name", config.ConfigName))

		if err := service.CreateConfig(ctx, transferConfig); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func validateConf(config server_dto.HTTPConfigDto) error {
	validate := validator.New()

	if err := validate.Struct(config); err != nil {
		return err
	}

	if err := config.Validate(); err != nil {
		return err
	}

	return nil
}

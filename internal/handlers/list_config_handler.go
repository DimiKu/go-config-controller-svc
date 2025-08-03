package handlers

import (
	"context"
	"go-config-controller-svc/dto/server_dto"
	"go-config-controller-svc/internal/interfaces"
	"go-config-controller-svc/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

func ListConfigHandler(service interfaces.ServerService, log *zap.Logger, ctx context.Context) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {

		configs, err := service.GetConfigsList(ctx)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		resConfigs := make([]server_dto.HTTPConfigDto, len(configs))

		for i, c := range configs {
			resConfigs[i].ConfigName = c.ConfigName
			resConfigs[i].ConfigBranch = c.ConfigBranch
			resConfigs[i].ConfigValue = c.ConfigValue
		}

		rw.Header().Set("Content-Type", "application/json")
		utils.MakeConfigsResponse(rw, resConfigs, log)
	}
}

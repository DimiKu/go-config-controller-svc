package utils

import (
	"encoding/json"
	"go-config-controller-svc/dto/server_dto"
	"go.uber.org/zap"
	"net/http"
)

func MakeConfigsResponse(w http.ResponseWriter, configs []server_dto.HTTPConfigDto, log *zap.Logger) {
	jsonRes, err := json.Marshal(configs)
	if err != nil {
		log.Error("can't decode response", zap.Error(err))
	}
	w.Write(jsonRes)
}

func MakeConfigResponse(w http.ResponseWriter, conf server_dto.HTTPConfigDto, log *zap.Logger) {
	jsonRes, err := json.Marshal(conf)
	if err != nil {
		log.Error("can't decode response", zap.Error(err))
	}
	w.Write(jsonRes)
}

func MakeTokenResponse(w http.ResponseWriter, token string, log *zap.Logger) {
	jsonRes, err := json.Marshal(token)
	if err != nil {
		log.Error("can't decode response", zap.Error(err))
	}
	w.Write(jsonRes)
}

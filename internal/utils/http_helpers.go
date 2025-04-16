package utils

import (
	"encoding/json"
	"go-config-controller-svc/dto/server_dto"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func MakeConfigsResponse(w http.ResponseWriter, configs []server_dto.HTTPConfigDto) {
	jsonRes, err := json.Marshal(configs)
	if err != nil {
		log.Fatal("can't decode response", zap.Error(err))
	}
	w.Write(jsonRes)
}

func MakeConfigResponse(w http.ResponseWriter, conf server_dto.HTTPConfigDto) {
	jsonRes, err := json.Marshal(conf)
	if err != nil {
		log.Fatal("can't decode response", zap.Error(err))
	}
	w.Write(jsonRes)
}

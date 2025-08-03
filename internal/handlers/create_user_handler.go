package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"go-config-controller-svc/dto/server_dto"
	"go-config-controller-svc/internal/entities"
	"go-config-controller-svc/internal/interfaces"
	"go-config-controller-svc/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

func CreateUserHandler(service interfaces.ServerService, log *zap.Logger, ctx context.Context) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		var user server_dto.UserDto
		var transferUser entities.User

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(buf.Bytes(), &user)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		hashPass := utils.GetSHA256Hash(user.Password)

		transferUser.Username = user.Username
		transferUser.Password = hashPass

		log.Info("Registration user:", zap.String("username", user.Username))

		if err := service.CreateUser(ctx, transferUser); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

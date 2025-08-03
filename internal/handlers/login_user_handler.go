package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"go-config-controller-svc/dto/server_dto"
	"go-config-controller-svc/internal/custom_errors"
	"go-config-controller-svc/internal/interfaces"
	"go-config-controller-svc/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

func LoginUserHandler(service interfaces.ServerService, log *zap.Logger, ctx context.Context) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		var user server_dto.UserDto

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

		appUser, err := service.GetUser(ctx, user.Username)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if appUser.Password != hashPass {
			http.Error(rw, custom_errors.ErrLoginError.Error(), http.StatusBadRequest)
			return
		}

		log.Info("Login user:", zap.String("username", user.Username))

		token, err := utils.CreateJWTToken(user.Username, []string{"/create_config", "/get_configs", "/delete_configs", "/execute"})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		utils.MakeTokenResponse(rw, token, log)
	}
}

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"go-config-controller-svc/dto/server_dto"
	"go-config-controller-svc/internal/interfaces"
	"go.uber.org/zap"
	"net/http"
)

func CreateTaskHandler(service interfaces.ServerService, log *zap.Logger, ctx context.Context) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		var task server_dto.TaskDto

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(buf.Bytes(), &task)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		log.Info("Got task:", zap.String("task", task.Task))

		if err := service.AddTaskForExecutor(ctx, task.Task); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

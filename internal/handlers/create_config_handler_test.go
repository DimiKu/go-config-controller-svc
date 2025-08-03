package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go-config-controller-svc/internal/entities"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MyConfRequest struct {
	ConfName   string `json:"config_name"`
	ConfValue  string `json:"config_value"`
	ConfBranch string `json:"config_branch"`
}

func TestCreateConfigHandler(t *testing.T) {
	log, _ := zap.NewProduction()
	mockCtl := gomock.NewController(t)
	mockSvc := NewMockServerService(mockCtl)

	type args struct {
		statusCode int
	}
	tests := []struct {
		name    string
		args    args
		req     MyConfRequest
		url     string
		wantErr bool
	}{
		{
			name: "positive test1",
			args: args{
				statusCode: http.StatusOK,
			},
			url:     "/create_config",
			wantErr: false,
			req: MyConfRequest{
				ConfName:   "nginx_values",
				ConfValue:  "test",
				ConfBranch: "main",
			},
		},
		{
			name: "neg test1",
			args: args{
				statusCode: http.StatusBadRequest,
			},
			url:     "/create_config",
			wantErr: false,
			req: MyConfRequest{
				ConfName:   "n",
				ConfValue:  "test",
				ConfBranch: "main",
			},
		},
		{
			name: "neg test2",
			args: args{
				statusCode: http.StatusBadRequest,
			},
			url:     "/create_config",
			wantErr: false,
			req: MyConfRequest{
				ConfName:   "nginx",
				ConfValue:  "t",
				ConfBranch: "main",
			},
		},
		{
			name: "neg test3",
			args: args{
				statusCode: http.StatusBadRequest,
			},
			url:     "/create_config",
			wantErr: false,
			req: MyConfRequest{
				ConfName:   ";drop database;",
				ConfValue:  "test",
				ConfBranch: "main",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var transferConfig entities.ServerConfig
			jsonData, _ := json.Marshal(tt.req)
			req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewBuffer(jsonData))
			w := httptest.NewRecorder()
			ctx := context.Background()
			handlerFunc := CreateConfigHandler(mockSvc, log, ctx)

			transferConfig.ConfigValue = tt.req.ConfValue
			transferConfig.ConfigName = tt.req.ConfName
			transferConfig.ConfigBranch = tt.req.ConfBranch

			mockSvc.EXPECT().CreateConfig(ctx, transferConfig).Return(nil).AnyTimes()

			handlerFunc(w, req)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.args.statusCode, result.StatusCode)
		})
	}

}

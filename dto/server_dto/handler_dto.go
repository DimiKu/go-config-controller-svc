package server_dto

import (
	"go-config-controller-svc/internal/custom_errors"
	"reflect"
	"strings"
)

type HTTPConfigDto struct {
	ConfigName   string `json:"config_name" validate:"required,min=3,max=20"`
	ConfigValue  string `json:"config_value" validate:"required,min=3,max=20"`
	ConfigBranch string `json:"config_branch" validate:"required,min=3,max=20"`
}

func (h *HTTPConfigDto) Validate() error {
	v := reflect.ValueOf(h).Elem() // Получаем значение структуры

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String {
			if strings.Contains(field.String(), ";") {
				return custom_errors.ErrFieldsContainsBadChars
			}
		}
	}
	return nil
}

type UserDto struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type TaskDto struct {
	Task string `json:"task" validate:"required,min=1,max=20"`
}

package server_dto

type HTTPConfigDto struct {
	ConfigName   string `json:"config_name"`
	ConfigValue  string `json:"config_value"`
	ConfigBranch string `json:"config_branch"`
}

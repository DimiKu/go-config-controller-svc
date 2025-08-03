package server_dto

type ConfigDBDto struct {
	ConfigName   string `db:"config_name"`
	ConfigValue  string `db:"config_value"`
	ConfigBranch string `db:"config_branch"`
}

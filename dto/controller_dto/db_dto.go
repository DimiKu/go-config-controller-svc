package controller_dto

import "time"

type ConfigDBDto struct {
	ConfigName   string    `db:"config_name"`
	ConfigValue  string    `db:"config_value"` // может и не надо
	LastUpdate   time.Time `db:"last_update"`
	ConfigBranch string    `db:"config_branch"`
	Locked       bool      `db:"locked"`
}

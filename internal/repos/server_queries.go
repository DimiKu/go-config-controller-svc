package repos

const (
	InsertNewConfig = `insert into configs (config_name, config_value, last_update, config_branch, locked) VALUES ($1, $2, NOW(), $3, false);`
)

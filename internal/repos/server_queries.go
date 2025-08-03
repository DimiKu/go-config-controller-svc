package repos

const (
	InsertNewConfig = `insert into configs (config_name, config_value, last_update, config_branch, locked) VALUES ($1, $2, NOW(), $3, false);`
	GetAllConfigs   = `SELECT config_name, config_value, config_branch from configs`
	GetCountConfigs = `SELECT count(*) from configs`
	DeleteConfig    = "DELETE FROM configs WHERE config_name = $1 and config_branch = $2"
)

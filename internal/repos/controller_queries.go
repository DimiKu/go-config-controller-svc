package repos

const (
	CreateConfigTable    = `CREATE TABLE IF NOT EXISTS configs (config_name varchar(255) NOT NULL, config_value TEXT, last_update TIMESTAMPTZ, config_branch text, locked bool);`
	GetRepoByName        = `SELECT config_name, config_value, last_update, config_branch FROM configs WHERE config_name = $1;`
	LockRepo             = `UPDATE configs SET locked = true WHERE config_name = $1;`
	UnlockRepo           = `UPDATE configs SET locked = false WHERE config_name = $1;`
	ChangeRepoUpdateTime = `UPDATE configs SET last_update = NOW() WHERE config_name = $1;`
	GetNotLockedConfig   = `SELECT config_name, config_value, last_update, config_branch, locked FROM configs WHERE locked = false ORDER BY last_update LIMIT 1`
	UnlockAllConfigs     = `UPDATE configs SET locked = false `
)

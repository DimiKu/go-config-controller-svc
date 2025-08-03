package entities

type ServerConfig struct {
	ConfigName   string
	ConfigValue  string
	ConfigBranch string
}

type User struct {
	Username string
	Password string
}

const (
	DoExecTask   = "DoExecTask"
	CommandQueue = "CommandQueue"
)

var (
	ExcludedPaths = map[string]bool{
		"/login":       true,
		"/create_user": true,
	}

	CheckPaths = []string{"/create_config", "/get_configs", "/delete_configs", "/execute"}
)

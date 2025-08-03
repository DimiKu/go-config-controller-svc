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

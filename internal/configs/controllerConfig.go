package configs

type ControllerConfig struct {
	GitRepo  string `env:"GIT_URL"`
	GitUser  string `env:"GIT_USER"`
	GitToken string `env:"GIT_TOKEN"`

	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBPort     string `env:"DB_PORT"`
	DBHost     string `env:"DB_HOST"`
	DBName     string `env:"DB_NAME"`
	Workers    int    `env:"WORKERS"`
}

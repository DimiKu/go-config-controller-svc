package configs

type ServerConfig struct {
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBPort     string `env:"DB_PORT"`
	DBHost     string `env:"DB_HOST"`
	DBName     string `env:"DB_NAME"`

	ServerAddr string `env:"SERVER_ADDR"`

	RedisAddr string `env:"REDIS_ADDR"`
	RedisDB   int    `env:"REDIS_DB"`
	Secret    string `env:"SECRET"`
}

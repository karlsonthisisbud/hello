package config

type Config struct {
	HTTPPort int
}

func New() *Config {
	return &Config{
		HTTPPort: 3000,
	}
}

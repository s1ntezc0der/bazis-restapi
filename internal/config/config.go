package config

import (
	"os"
	"strconv"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	Secret string
	Expire int
}

type Config struct {
	DB     DBConfig
	Server ServerConfig
	JWT    JWTConfig
}

func Load() Config {
	return Config{
		DB: DBConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "supersecret"),
			Expire: getEnvInt("JWT_EXPIRE", 24),
		},
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

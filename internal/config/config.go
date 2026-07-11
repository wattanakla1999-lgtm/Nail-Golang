package config

import "os"

type Config struct {
	Port string
	DSN  string
}

func Load() Config {
	return Config{
		Port: getEnv("PORT", "8080"),
		DSN:  getEnv("DATABASE_DSN", "host=localhost user=nailly password=nailly1234 dbname=nailly_db port=5432 sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

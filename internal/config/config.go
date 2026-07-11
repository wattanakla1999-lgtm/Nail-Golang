package config

import (
	"os"
	"strings"
)

type Config struct {
	Port string
	DSN  string
}

func Load() Config {
	loadEnvFile(".env")

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

func loadEnvFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

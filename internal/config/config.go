package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	TemplatesDir string
	GotenbergURL string
}

func Load() *Config {
	_ = godotenv.Load() // Load .env file if it exists (ignores error if missing)

	return &Config{
		Port:         getEnv("PORT", "8080"),
		TemplatesDir: getEnv("TEMPLATES_DIR", "./templates"),
		GotenbergURL: getEnv("GOTENBERG_URL", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

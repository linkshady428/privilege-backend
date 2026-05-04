package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	SkipAuth    bool
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		Port:        port,
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		SkipAuth:    os.Getenv("SKIP_AUTH") == "true",
	}
}

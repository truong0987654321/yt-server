package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTAccessSecret     string
	JWTRefreshSecret    string
	JWTAccessTTLMinutes int
	JWTRefreshTTLDays   int

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	FrontendCallbackURL string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	return &Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: mustGetEnv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME", "yt_clone"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		JWTAccessSecret:     mustGetEnv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret:    mustGetEnv("JWT_REFRESH_SECRET"),
		JWTAccessTTLMinutes: getEnvAsInt("JWT_ACCESS_TTL_MINUTES", 15),
		JWTRefreshTTLDays:   getEnvAsInt("JWT_REFRESH_TTL_DAYS", 30),

		GoogleClientID:     mustGetEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: mustGetEnv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  mustGetEnv("GOOGLE_REDIRECT_URL"),

		FrontendCallbackURL: getEnv("FRONTEND_CALLBACK_URL", "http://localhost:3000/auth/callback"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		log.Fatalf("Required environment variable is missing: %s", key)
	}
	return v
}

func getEnvAsInt(key string, fallback int) int {
	v := getEnv(key, "")
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}

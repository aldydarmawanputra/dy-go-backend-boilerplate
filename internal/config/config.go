package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppHost string
	AppPort string
	AppEnv  string

	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	JWTSecret      string
	JWTExpireHours int

	CORSAllowOrigins     string
	CORSAllowMethods     string
	CORSAllowHeaders     string
	CORSAllowCredentials bool

	AutoMigrate bool
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("config: no .env file found, relying on environment variables")
	}

	return &Config{
		AppHost: getEnv("APP_HOST", "0.0.0.0"),
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvInt("DB_PORT", 5432),
		DBName:     getEnv("DB_NAME", "aldy_dev"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnvInt("REDIS_PORT", 6379),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 24),

		CORSAllowOrigins:     getEnv("CORS_ALLOW_ORIGINS", "*"),
		CORSAllowMethods:     getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"),
		CORSAllowHeaders:     getEnv("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization"),
		CORSAllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", false),

		AutoMigrate: getEnvBool("AUTO_MIGRATE", false),
	}
}

func (c *Config) IsProduction() bool { return c.AppEnv == "production" }

func (c *Config) Addr() string { return c.AppHost + ":" + c.AppPort }

func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(c.DBUser),
		url.QueryEscape(c.DBPassword),
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

func (c *Config) RedisAddr() string { return fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort) }

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

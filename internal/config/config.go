package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

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

	AppBaseURL string

	JWTSecret              string
	JWTIssuer              string
	JWTAudience            string
	JWTExpireHours         int
	RefreshExpireHours     int
	VerifyTokenExpireHours int
	ResetTokenExpireHours  int

	CORSAllowOrigins     string
	CORSAllowMethods     string
	CORSAllowHeaders     string
	CORSAllowCredentials bool

	RateLimitMax       int
	RateLimitWindowSec int
	AuthRateLimitMax   int

	ReadTimeoutSec  int
	WriteTimeoutSec int
	IdleTimeoutSec  int
	BodyLimitBytes  int

	OTelEnabled     bool
	OTelServiceName string
	OTelEndpoint    string
	OTelInsecure    bool

	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	PaymentProvider      string
	PaymentCurrency      string
	PaymentWebhookSecret string

	StorageDriver        string
	StorageLocalPath     string
	StoragePublicBaseURL string

	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2Bucket          string
	R2PublicBaseURL   string

	SupabaseURL           string
	SupabaseServiceKey    string
	SupabaseBucket        string
	SupabasePublicBaseURL string

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

		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		JWTIssuer:   getEnv("JWT_ISSUER", "go-backend-boilerplate"),
		JWTAudience: getEnv("JWT_AUDIENCE", "go-backend-boilerplate"),
		AppBaseURL:  getEnv("APP_BASE_URL", "http://localhost:8080"),

		JWTExpireHours:         getEnvInt("JWT_EXPIRE_HOURS", 24),
		RefreshExpireHours:     getEnvInt("REFRESH_TOKEN_EXPIRE_HOURS", 168),
		VerifyTokenExpireHours: getEnvInt("VERIFY_TOKEN_EXPIRE_HOURS", 24),
		ResetTokenExpireHours:  getEnvInt("RESET_TOKEN_EXPIRE_HOURS", 1),

		CORSAllowOrigins:     getEnv("CORS_ALLOW_ORIGINS", "*"),
		CORSAllowMethods:     getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"),
		CORSAllowHeaders:     getEnv("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization"),
		CORSAllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", false),

		RateLimitMax:       getEnvInt("RATE_LIMIT_MAX", 100),
		RateLimitWindowSec: getEnvInt("RATE_LIMIT_WINDOW_SEC", 60),
		AuthRateLimitMax:   getEnvInt("AUTH_RATE_LIMIT_MAX", 10),

		ReadTimeoutSec:  getEnvInt("READ_TIMEOUT_SEC", 10),
		WriteTimeoutSec: getEnvInt("WRITE_TIMEOUT_SEC", 15),
		IdleTimeoutSec:  getEnvInt("IDLE_TIMEOUT_SEC", 60),
		BodyLimitBytes:  getEnvInt("BODY_LIMIT_BYTES", 1048576),

		OTelEnabled:     getEnvBool("OTEL_ENABLED", false),
		OTelServiceName: getEnv("OTEL_SERVICE_NAME", "go-backend-boilerplate"),
		OTelEndpoint:    getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318"),
		OTelInsecure:    getEnvBool("OTEL_EXPORTER_OTLP_INSECURE", true),

		SMTPHost:     getEnv("SMTP_HOST", "localhost"),
		SMTPPort:     getEnv("SMTP_PORT", "1025"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "no-reply@example.com"),

		PaymentProvider:      getEnv("PAYMENT_PROVIDER", "stub"),
		PaymentCurrency:      getEnv("PAYMENT_CURRENCY", "IDR"),
		PaymentWebhookSecret: getEnv("PAYMENT_WEBHOOK_SECRET", ""),

		StorageDriver:        getEnv("STORAGE_DRIVER", "local"),
		StorageLocalPath:     getEnv("STORAGE_LOCAL_PATH", "./storage"),
		StoragePublicBaseURL: getEnv("STORAGE_PUBLIC_BASE_URL", "http://localhost:8080/storage"),

		R2AccountID:       getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2Bucket:          getEnv("R2_BUCKET", ""),
		R2PublicBaseURL:   getEnv("R2_PUBLIC_BASE_URL", ""),

		SupabaseURL:           getEnv("SUPABASE_URL", ""),
		SupabaseServiceKey:    getEnv("SUPABASE_SERVICE_KEY", ""),
		SupabaseBucket:        getEnv("SUPABASE_BUCKET", ""),
		SupabasePublicBaseURL: getEnv("SUPABASE_PUBLIC_BASE_URL", ""),

		AutoMigrate: getEnvBool("AUTO_MIGRATE", false),
	}
}

func (c *Config) IsProduction() bool { return c.AppEnv == "production" }

// Validate fails fast on dangerous misconfiguration (mostly in production).
func (c *Config) Validate() error {
	var problems []string
	if c.JWTSecret == "" {
		problems = append(problems, "JWT_SECRET is empty")
	}
	if c.IsProduction() {
		if c.JWTSecret == "change-me-in-production" || len(c.JWTSecret) < 16 {
			problems = append(problems, "JWT_SECRET must be a strong value (>=16 chars) in production")
		}
		if c.CORSAllowCredentials && c.CORSAllowOrigins == "*" {
			problems = append(problems, "CORS_ALLOW_CREDENTIALS=true is not allowed with CORS_ALLOW_ORIGINS=*")
		}
	}
	if len(problems) > 0 {
		return fmt.Errorf("invalid configuration: %s", strings.Join(problems, "; "))
	}
	return nil
}

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

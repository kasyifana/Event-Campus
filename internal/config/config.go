package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	Supabase   SupabaseConfig
	PostgreSQL PostgreSQLConfig
	JWT        JWTConfig
	Email      EmailConfig
	Upload     UploadConfig
	CORS       CORSConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type SupabaseConfig struct {
	URL        string
	AnonKey    string
	ServiceKey string
}

type PostgreSQLConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	Expiration string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
}

type UploadConfig struct {
	MaxSize int64
	Path    string
}

type CORSConfig struct {
	AllowedOrigins []string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (ignore error in production)
	_ = godotenv.Load()

	smtpPort, err := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}

	maxSize, err := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_UPLOAD_SIZE: %w", err)
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Supabase: SupabaseConfig{
			URL:        getEnvRequired("SUPABASE_URL"),
			AnonKey:    getEnvRequired("SUPABASE_ANON_KEY"),
			ServiceKey: getEnvRequired("SUPABASE_SERVICE_KEY"),
		},
		PostgreSQL: PostgreSQLConfig{
			Host:     getEnv("POSTGRES_HOST", "aws-0-ap-southeast-1.pooler.supabase.com"),
			Port:     getEnv("POSTGRES_PORT", "6543"),
			Database: getEnv("POSTGRES_DB", "postgres"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnvRequired("POSTGRES_PASSWORD"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "require"),
		},
		JWT: JWTConfig{
			Secret:     getEnvRequired("JWT_SECRET"),
			Expiration: getEnv("JWT_EXPIRATION", "24h"),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     smtpPort,
			SMTPUser:     getEnvRequired("SMTP_USER"),
			SMTPPassword: getEnvRequired("SMTP_PASSWORD"),
		},
		Upload: UploadConfig{
			MaxSize: maxSize,
			Path:    getEnv("UPLOAD_PATH", "./storage"),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		},
	}

	return config, nil
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvRequired gets required environment variable and panics if not set
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

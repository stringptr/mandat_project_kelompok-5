package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Host               string
	Port               string
	CORSOrigins        []string
	DBMasterConfig     PostgresConfig
	AuthConfig         AuthConfig
	NATSConfig         NATSConfig
	RestrictAuthConfig RestrictAuthConfig
	MailConfig         MailConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type AuthConfig struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type NATSConfig struct {
	Host  string
	Port  string
	Token string
}

type RestrictAuthConfig struct {
	MaxAttempt int
	Duration   time.Duration
}

type MailConfig struct {
	Provider     string
	Host         string
	Port         int
	Username     string
	Password     string
	FromEmail    string
	FromName     string
	ResendAPIKey string
}

func Load() *Config {
	serverHost := getEnv("BACKEND_IP", "0.0.0.0")
	serverPort := getEnv("BACKEND_PORT", "8080")
	frontendHost := getEnv("FRONTEND_IP", "localhost")
	frontendPort := getEnv("FRONTEND_PORT", "5173")

	dbHost := getEnv("MASTER_HOST", "master")
	dbName := getEnv("MASTER_DB", "imunisasi")
	dbUser := getEnv("MASTER_USER", "postgres")
	dbPass := getEnv("MASTER_PASSWORD", "postgres")

	jwtSecret := getEnv("JWT_SECRET", "abcdefghijklmnopqrstuvwxyz")

	natsHost := getEnv("NATS_HOST", "nats")
	natsPort := getEnv("NATS_PORT", "4222")
	natsToken := getEnv("NATS_TOKEN", "")

	restrictAuthDuration := getEnv("BANNED_AUTH_DURATION", "15")
	restrictAuthMaxAttempt := getEnv("BANNED_AUTH_MAX_ATTEMPT", "3")
	restrictAuthDurationInt, _ := strconv.Atoi(restrictAuthDuration)
	restrictAuthMaxAttemptInt, _ := strconv.Atoi(restrictAuthMaxAttempt)

	mailProvider := getEnv("MAIL_PROVIDER", "smtp")
	mailHost := getEnv("MAIL_HOST", "sandbox.smtp.mailtrap.io")
	mailPortStr := getEnv("MAIL_PORT", "2525")
	mailUser := getEnv("MAIL_USERNAME", "")
	mailPass := getEnv("MAIL_PASSWORD", "")
	mailFromEmail := getEnv("MAIL_FROM_EMAIL", "noreply@sigizi.com")
	mailFromName := getEnv("MAIL_FROM_NAME", "SiGizi")
	mailPortInt, _ := strconv.Atoi(mailPortStr)
	resendAPIKey := getEnv("RESEND_API_KEY", "")

	return &Config{
		Host:        serverHost,
		Port:        serverPort,
		CORSOrigins: []string{fmt.Sprintf("http://%s:%s", frontendHost, frontendPort)},
		DBMasterConfig: PostgresConfig{
			Host:     dbHost,
			Port:     "5432",
			User:     dbUser,
			Password: dbPass,
			DBName:   dbName,
			SSLMode:  "disable",
		},
		AuthConfig: AuthConfig{
			JWTSecret:       jwtSecret,
			AccessTokenTTL:  30 * time.Minute,
			RefreshTokenTTL: 7 * 24 * time.Hour,
		},
		NATSConfig: NATSConfig{
			Host:  natsHost,
			Port:  natsPort,
			Token: natsToken,
		},
		RestrictAuthConfig: RestrictAuthConfig{
			MaxAttempt: restrictAuthMaxAttemptInt,
			Duration:   time.Duration(restrictAuthDurationInt) * time.Second,
		},
		MailConfig: MailConfig{
			Provider:     mailProvider,
			Host:         mailHost,
			Port:         mailPortInt,
			Username:     mailUser,
			Password:     mailPass,
			FromEmail:    mailFromEmail,
			FromName:     mailFromName,
			ResendAPIKey: resendAPIKey,
		},
	}
}

func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

func (c *NATSConfig) URL() string {
	return fmt.Sprintf(
		"nats://%s:%s",
		c.Host, c.Port)
}

func (c *Config) Server() (string, string) {
	return c.Host, c.Port
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

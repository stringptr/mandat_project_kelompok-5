package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Host           string
	Port           string
	CORSOrigins    []string
	DBMasterConfig PostgresConfig
	NATSConfig     NATSConfig
	AuthConfig     AuthConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type NATSConfig struct {
	natsHost string
	natsPort string
}

type AuthConfig struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
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

	natsHost := getEnv("NATS_HOST", "message-system")
	natsPort := getEnv("NATS_PORT", "4222")

	jwtSecret := getEnv("JWT_SECRET", "abcdefghijklmnopqrstuvwxyz")

	return &Config{
		Host:        serverHost,
		Port:        serverPort,
		CORSOrigins: []string{fmt.Sprintf("%q:%q", frontendHost, frontendPort)},
		DBMasterConfig: PostgresConfig{
			Host:     dbHost,
			Port:     "5432",
			User:     dbUser,
			Password: dbPass,
			DBName:   dbName,
			SSLMode:  "disable",
		},
		NATSConfig: NATSConfig{
			natsHost: natsHost,
			natsPort: natsPort,
		},
		AuthConfig: AuthConfig{
			JWTSecret:       jwtSecret,
			AccessTokenTTL:  30 * time.Minute,
			RefreshTokenTTL: 7 * 24 * time.Hour,
		},
	}
}

func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

func (c *NATSConfig) NATSURL() string {
	return fmt.Sprintf("nats://%s:%s", c.natsHost, c.natsPort)
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

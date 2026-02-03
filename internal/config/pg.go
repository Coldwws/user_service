package config

import "fmt"

type PGConfig struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
}

func loadPG() PGConfig {
	return PGConfig{
		Host:     getEnv("PG_HOST", "localhost"),
		Port:     getEnv("PG_PORT", "5432"),
		DBName:   getEnv("PG_DATABASE_NAME", "user_service"),
		User:     getEnv("PG_USER", "postgres"),
		Password: getEnv("PG_PASSWORD", "postgres"),
		SSLMode:  getEnv("PG_SSLMODE", "disable"),
	}
}

func (c PGConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.DBName, c.User, c.Password, c.SSLMode,
	)
}

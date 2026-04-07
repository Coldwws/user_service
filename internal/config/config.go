package config

import (
	"log"
	"os"
)

type Config struct {
	Env     string
	GRPC    GRPCConfig
	PG      PGConfig
	Http    HttpConfig
	Swagger SwaggerConfig
}

func LoadConfig() Config {
	httpConf, err := NewHttpConfig()
	if err != nil {
		log.Fatalf("failed to load HTTP config: %v", err)
	}

	swaggerConf, err := NewSwaggerConfig()
	if err != nil {
		log.Fatalf("failed to load Swagger config: %v", err)
	}

	return Config{
		Env:     getEnv("ENV", "local"),
		GRPC:    loadGRPC(),
		PG:      loadPG(),
		Http:    httpConf,
		Swagger: swaggerConf,
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing env var :%s", key)
	}
	return v
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	return v
}

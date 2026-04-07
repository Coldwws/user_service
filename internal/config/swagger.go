package config

import (
	"github.com/pkg/errors"
	"net"
	"os"
)

type SwaggerConfig interface {
	Address() string
}

type swaggerConfig struct {
	Host string
	Port string
}

func NewSwaggerConfig() (SwaggerConfig, error) {
	host := os.Getenv("SWAGGER_HOST")
	if len(host) == 0 {
		return nil, errors.New("swagger host not found")
	}
	port := os.Getenv("SWAGGER_PORT")
	if len(port) == 0 {
		return nil, errors.New("swagger port not found")
	}
	return &swaggerConfig{
		Host: host,
		Port: port,
	}, nil
}

func (s *swaggerConfig) Address() string {
	return net.JoinHostPort(s.Host, s.Port)
}

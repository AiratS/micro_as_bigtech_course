package config

import (
	"fmt"
	"net"
	"os"
)

const (
	SwaggerHostEnvName = "SWAGGER_HOST"
	SwaggerPortEnvName = "SWAGGER_PORT"
)

type SwaggerConfig interface {
	Address() string
}

type swaggerConfig struct {
	host string
	port string
}

func NewSwaggerConfig() (SwaggerConfig, error) {
	host := os.Getenv(SwaggerHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("SWAGGER_HOST is not found")
	}

	port := os.Getenv(SwaggerPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("SWAGGER_PORT is not found")
	}

	return &swaggerConfig{
		host: host,
		port: port,
	}, nil
}

func (c *swaggerConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}

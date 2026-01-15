package config

import (
	"errors"
	"net"
	"os"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type GRPCConfig interface {
	Address() string
}

type grpcConfig struct {
	host string
	port string
}

func (c *grpcConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}

func NewGRPCConfig() (GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host is not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port is not found")
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

// Package config provides server configuration
package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

// Server represents server configuration
type Server struct {
	Address     string `env:"ADDRESS" envDefault:"localhost:6060"`
	DatabaseDSN string `env:"DATABASE_DSN,required"`
}

// NewServer returns new server configuration object
func NewServer() *Server {
	return &Server{}
}

// ParseEnvs reads server configuration from environment
func ParseEnvs(c *Server) error {
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("failed to parse server config: %w", err)
	}
	return nil
}

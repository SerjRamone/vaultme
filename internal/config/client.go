package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Client struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:6161"`
}

func NewClient() *Client {
	return &Client{}
}

// ParseClientEnvs reads client configuration from environment
func ParseClientEnvs(c *Client) error {
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("failed to parse client config: %w", err)
	}
	return nil
}

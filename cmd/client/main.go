// Package main provides entrypoint for client
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/SerjRamone/vaultme/internal/client"
	"github.com/SerjRamone/vaultme/internal/config"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("an error occured: %v", err)
	}
}

func run() error {
	ctx := context.Background()

	cfg := config.NewClient()
	if err := config.ParseClientEnvs(cfg); err != nil {
		return fmt.Errorf("parse config envs error: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	app := client.NewApp(logger)

	if err := app.Run(ctx, cfg); err != nil {
		return fmt.Errorf("failed to run app: %w", err)
	}

	return nil
}

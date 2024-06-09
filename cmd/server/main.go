// Package main provides entrypoint for server
package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/SerjRamone/vaultme/internal/config"
	"github.com/SerjRamone/vaultme/internal/repository"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("an error occured: %v", err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to create logger: %w ", err)
	}

	// TODO check error
	defer logger.Sync()

	cfg := config.NewServer()
	if err := config.ParseEnvs(cfg); err != nil {
		return fmt.Errorf("failed to parse server config: %w", err)
	}

	db, err := repository.NewDB(ctx, cfg.DatabaseDSN, logger)
	if err != nil {
		return fmt.Errorf("failed to create db: %w", err)
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer logger.Info("closing DB")
		defer wg.Done()
		<-ctx.Done()

		db.Close()
	}()

	logger.Info("DB is running")

	// TODO run server

	logger.Info("server is running")

	defer func() {
		wg.Wait()
	}()

	select {
	case <-ctx.Done():
		// TODO
	}

	go func() {
		ctx, cancelTOut := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelTOut()

		<-ctx.Done()
		logger.Error("shutdown timeout")
	}()
	return nil
}

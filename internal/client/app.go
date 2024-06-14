// Package client ...
package client

import (
	"context"

	"github.com/SerjRamone/vaultme/internal/config"
	"go.uber.org/zap"
)

// App - is a VaultMe client app
type App struct {
	ui  *TerminalUI
	log *zap.Logger
}

// NewApp - returns new App instance
func NewApp(log *zap.Logger) *App {
	return &App{
		ui:  &TerminalUI{},
		log: log,
	}
}

// Run - starts VaultMe app
func (a *App) Run(ctx context.Context, cfg *config.Client) error {
	return a.ui.Start(ctx, a.log, cfg)
}

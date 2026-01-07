package main

import (
	"context"
	"log/slog"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"MMDContent/internal/services/fs"
	"MMDContent/internal/storage"
)

type App struct {
	ctx           context.Context
	modelsStorage *storage.Models
	stagesStorage *storage.Stages
}

func NewApp(modelsStorage *storage.Models, stagesStorage *storage.Stages) *App {
	return &App{
		modelsStorage: modelsStorage,
		stagesStorage: stagesStorage,
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	modelsData, stagesData, err := fs.SyncFiles()
	if err != nil {
		slog.InfoContext(ctx, "error in startup", "error", err)
	}

	a.modelsStorage.Set(modelsData)
	a.stagesStorage.Set(stagesData)
}

// Quit closes the app
func (a *App) Quit() {
	wailsruntime.Quit(a.ctx)
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {}

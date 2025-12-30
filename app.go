package main

import (
	"context"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx         context.Context
	isMinimized bool
}

func NewApp() *App {
	return &App{
		isMinimized: false,
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Quit cierra la aplicación
func (a *App) Quit() {
	wailsruntime.Quit(a.ctx)
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {

}

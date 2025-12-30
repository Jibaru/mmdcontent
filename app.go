package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx         context.Context
	isMinimized bool
	modelsData  *ModelsData
}

func NewApp() *App {
	return &App{
		isMinimized: false,
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Parse models data on startup
	err := ParseModelsData()
	if err != nil {
		fmt.Printf("Error parsing models data: %v\n", err)
	} else {
		fmt.Println("Models data parsed successfully")
	}

	// Load models data
	modelsData, err := LoadModelsData()
	if err != nil {
		fmt.Printf("Error loading models data: %v\n", err)
	} else {
		a.modelsData = modelsData
		fmt.Printf("Loaded %d models\n", len(modelsData.Models))
	}
}

// GetModels returns paginated models
func (a *App) GetModels(page, perPage int) PaginatedModels {
	if a.modelsData == nil {
		return PaginatedModels{
			Models:     []Model{},
			Total:      0,
			Page:       page,
			PerPage:    perPage,
			TotalPages: 0,
		}
	}

	return a.modelsData.GetPaginatedModels(page, perPage)
}

// GetAllModels returns all models without pagination
func (a *App) GetAllModels() []Model {
	if a.modelsData == nil {
		return []Model{}
	}

	return a.modelsData.Models
}

// RefreshModelsData re-parses and reloads the models data
func (a *App) RefreshModelsData() error {
	err := ParseModelsData()
	if err != nil {
		return err
	}

	modelsData, err := LoadModelsData()
	if err != nil {
		return err
	}

	a.modelsData = modelsData
	return nil
}

// GetImageAsBase64 reads an image file and returns it as base64 string
func (a *App) GetImageAsBase64(filePath string) (string, error) {
	imageData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Determine mime type based on file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	mimeType := "image/jpeg"
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	}

	base64Data := base64.StdEncoding.EncodeToString(imageData)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data), nil
}

// Quit cierra la aplicación
func (a *App) Quit() {
	wailsruntime.Quit(a.ctx)
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {

}

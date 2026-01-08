package main

import (
	"embed"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"MMDContent/internal/handlers"
	"MMDContent/internal/services/openai"
	"MMDContent/internal/storage"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed icon.ico
var icon []byte

func main() {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	modelsStorage, err := storage.NewModelsLoaded(filepath.Join("data", "Models"), filepath.Join("data", "models.json"))
	if err != nil {
		slog.Error("error loading models", "error", err)
		return
	}

	stagesStorage, err := storage.NewStagesLoaded(filepath.Join("data", "Stages"), filepath.Join("data", "stages.json"))
	if err != nil {
		slog.Error("error loading stages", "error", err)
		return
	}

	motionsStorage, err := storage.NewMotionsLoaded(filepath.Join("data", "Motions"), filepath.Join("data", "motions.json"))
	if err != nil {
		slog.Error("error loading motions", "error", err)
		return
	}

	images := handlers.NewImages()
	embeddings := handlers.NewEmbeddings(*client)
	models := handlers.NewModels(*client, modelsStorage)
	stages := handlers.NewStages(*client, stagesStorage)
	motions := handlers.NewMotions(*client, motionsStorage)

	app := NewApp(modelsStorage, stagesStorage)

	err = wails.Run(&options.App{
		Title:            "MMDContent",
		Width:            1080,
		Height:           720,
		WindowStartState: options.Maximised,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
			images,
			embeddings,
			models,
			stages,
			motions,
		},
	})

	if err != nil {
		log.Fatal("Error starting app:", err)
	}
}

package main

import (
	"embed"
	"log"
	"os"

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
	modelsStorage := storage.NewModels()
	stagesStorage := storage.NewStages()

	images := handlers.NewImages()
	embeddings := handlers.NewEmbeddings(*client)
	models := handlers.NewModels(*client, modelsStorage)
	stages := handlers.NewStages(*client, stagesStorage)

	app := NewApp(modelsStorage, stagesStorage)

	err := wails.Run(&options.App{
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
		},
	})

	if err != nil {
		log.Fatal("Error starting app:", err)
	}
}

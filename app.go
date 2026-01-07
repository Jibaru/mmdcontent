package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx         context.Context
	isMinimized bool
	modelsData  *ModelsData
	stagesData  *StagesData
}

func NewApp() *App {
	return &App{
		isMinimized: false,
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	fmt.Println("\n========================================")
	fmt.Println("🚀 Starting Application")
	fmt.Println("========================================\n")

	// Check for models changes
	fmt.Println("📦 Checking Models directory...")
	modelsChanges, err := CheckModelsChanges()
	if err != nil {
		fmt.Printf("⚠️  Error checking models: %v\n", err)
	}
	PrintChanges("Models", modelsChanges)

	// Parse models data if changes detected or JSON doesn't exist
	if NeedsSync(modelsChanges) || !fileExists("data/data.json") {
		fmt.Println("   🔄 Syncing models data...")
		err := ParseModelsData()
		if err != nil {
			fmt.Printf("   ❌ Error parsing models data: %v\n", err)
		} else {
			fmt.Println("   ✅ Models data synced successfully")
		}
	}

	// Load models data
	modelsData, err := LoadModelsData()
	if err != nil {
		fmt.Printf("❌ Error loading models data: %v\n", err)
	} else {
		a.modelsData = modelsData
		fmt.Printf("✅ Loaded %d models\n\n", len(modelsData.Models))
	}

	// Check for stages changes
	fmt.Println("🎭 Checking Stages directory...")
	stagesChanges, err := CheckStagesChanges()
	if err != nil {
		fmt.Printf("⚠️  Error checking stages: %v\n", err)
	}
	PrintChanges("Stages", stagesChanges)

	// Parse stages data if changes detected or JSON doesn't exist
	if NeedsSync(stagesChanges) || !fileExists("data/stages.json") {
		fmt.Println("   🔄 Syncing stages data...")
		err = ParseStagesData()
		if err != nil {
			fmt.Printf("   ❌ Error parsing stages data: %v\n", err)
		} else {
			fmt.Println("   ✅ Stages data synced successfully")
		}
	}

	// Load stages data
	stagesData, err := LoadStagesData()
	if err != nil {
		fmt.Printf("❌ Error loading stages data: %v\n", err)
	} else {
		a.stagesData = stagesData
		fmt.Printf("✅ Loaded %d stages\n", len(stagesData.Stages))
	}

	fmt.Println("\n========================================")
	fmt.Println("✅ Application Ready")
	fmt.Println("========================================\n")
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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

// GetStages returns paginated stages
func (a *App) GetStages(page, perPage int) PaginatedStages {
	if a.stagesData == nil {
		return PaginatedStages{
			Stages:     []Stage{},
			Total:      0,
			Page:       page,
			PerPage:    perPage,
			TotalPages: 0,
		}
	}

	return a.stagesData.GetPaginatedStages(page, perPage)
}

// GetAllStages returns all stages without pagination
func (a *App) GetAllStages() []Stage {
	if a.stagesData == nil {
		return []Stage{}
	}

	return a.stagesData.Stages
}

// RefreshStagesData re-parses and reloads the stages data
func (a *App) RefreshStagesData() error {
	err := ParseStagesData()
	if err != nil {
		return err
	}

	stagesData, err := LoadStagesData()
	if err != nil {
		return err
	}

	a.stagesData = stagesData
	return nil
}

// SearchModels searches models using semantic similarity with embeddings
func (a *App) SearchModels(query string, limit int) ([]Model, error) {
	if a.modelsData == nil {
		return []Model{}, nil
	}

	// Generate embedding for the search query
	queryEmbedding, err := GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Calculate similarity scores for all models
	type scoredModel struct {
		model Model
		score float64
	}

	var scoredModels []scoredModel
	for _, model := range a.modelsData.Models {
		if len(model.Embedding) == 0 {
			// Skip models without embeddings
			continue
		}

		similarity := CosineSimilarity(queryEmbedding, model.Embedding)
		scoredModels = append(scoredModels, scoredModel{
			model: model,
			score: similarity,
		})
	}

	// Sort by similarity score (descending)
	sort.Slice(scoredModels, func(i, j int) bool {
		return scoredModels[i].score > scoredModels[j].score
	})

	// Return top results
	if limit <= 0 || limit > len(scoredModels) {
		limit = len(scoredModels)
	}

	results := make([]Model, limit)
	for i := 0; i < limit; i++ {
		results[i] = scoredModels[i].model
	}

	return results, nil
}

// SearchStages searches stages using semantic similarity with embeddings
func (a *App) SearchStages(query string, limit int) ([]Stage, error) {
	if a.stagesData == nil {
		return []Stage{}, nil
	}

	// Generate embedding for the search query
	queryEmbedding, err := GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Calculate similarity scores for all stages
	type scoredStage struct {
		stage Stage
		score float64
	}

	var scoredStages []scoredStage
	for _, stage := range a.stagesData.Stages {
		if len(stage.Embedding) == 0 {
			// Skip stages without embeddings
			continue
		}

		similarity := CosineSimilarity(queryEmbedding, stage.Embedding)
		scoredStages = append(scoredStages, scoredStage{
			stage: stage,
			score: similarity,
		})
	}

	// Sort by similarity score (descending)
	sort.Slice(scoredStages, func(i, j int) bool {
		return scoredStages[i].score > scoredStages[j].score
	})

	// Return top results
	if limit <= 0 || limit > len(scoredStages) {
		limit = len(scoredStages)
	}

	results := make([]Stage, limit)
	for i := 0; i < limit; i++ {
		results[i] = scoredStages[i].stage
	}

	return results, nil
}

// GenerateEmbeddingsForAll generates embeddings for all models and stages
func (a *App) GenerateEmbeddingsForAll() error {
	return GenerateAllEmbeddings()
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

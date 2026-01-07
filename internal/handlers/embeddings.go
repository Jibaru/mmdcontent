package handlers

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"MMDContent/internal/entities"
	"MMDContent/internal/services/openai"
	"MMDContent/internal/storage"
)

type Embeddings struct {
	client        openai.Client
	modelsStorage *storage.Models
	stagesStorage *storage.Stages
}

func NewEmbeddings(
	client openai.Client,
) *Embeddings {
	return &Embeddings{
		client: client,
	}
}

func (e *Embeddings) GenerateAll() {
	// Generate embeddings for models
	if err := e.GenerateModelsEmbeddings(); err != nil {
		slog.Error("error generating model embeddings", "error", err)
		return
	}

	// Generate embeddings for stages
	if err := e.GenerateStagesEmbeddings(); err != nil {
		slog.Error("error generating stage embeddings", "error", err)
		return
	}
}

// GenerateModelsEmbeddings generates embeddings for all models
func (e *Embeddings) GenerateModelsEmbeddings() error {
	err := e.modelsStorage.Refresh()
	if err != nil {
		return err
	}

	totalModels := e.modelsStorage.Total()
	skippedCount := 0
	updatedCount := 0
	failedCount := 0

	fmt.Printf("   Found %d models total\n", totalModels)

	updatedModels := make([]entities.Model, totalModels)
	for i, model := range e.modelsStorage.Get().Models {
		// Skip if embedding already exists
		if len(model.Embedding) > 0 {
			skippedCount++
			continue
		}

		// Prepare text for embedding
		text := PrepareTextForEmbedding(model.Name, model.Description)

		fmt.Printf("   [%d/%d] Generating embedding for: %s\n", i+1, totalModels, model.Name)

		// Generate embedding
		embedding, err := e.client.GenerateEmbedding(text)
		if err != nil {
			fmt.Printf("   ⚠️  Warning: Failed for %s: %v\n", model.ID, err)
			failedCount++
			continue
		}

		model.Embedding = embedding
		updatedCount++

		// Small delay to avoid rate limits (optional, adjust based on your API tier)
		time.Sleep(100 * time.Millisecond)

		updatedModels[i] = model
	}

	fmt.Printf("\n   ✅ Generated: %d | ⏭️  Skipped: %d | ❌ Failed: %d\n", updatedCount, skippedCount, failedCount)

	// Save updated data back to file
	if updatedCount > 0 {
		fmt.Println("   💾 Saving models data...")
		e.modelsStorage.Set(&entities.ModelsData{Models: updatedModels})
		if err := e.modelsStorage.Save(); err != nil {
			return err
		}
		fmt.Println("   ✅ Models data saved successfully")
	} else {
		fmt.Println("   ℹ️  No new embeddings to save")
	}

	return nil
}

// GenerateStagesEmbeddings generates embeddings for all stages
func (e *Embeddings) GenerateStagesEmbeddings() error {
	err := e.modelsStorage.Refresh()
	if err != nil {
		return err
	}

	totalStages := e.stagesStorage.Total()
	skippedCount := 0
	updatedCount := 0
	failedCount := 0

	fmt.Printf("   Found %d stages total\n", totalStages)

	updatedStages := make([]entities.Stage, totalStages)
	for i, stage := range e.stagesStorage.Get().Stages {
		// Skip if embedding already exists
		if len(stage.Embedding) > 0 {
			skippedCount++
			continue
		}

		// Prepare text for embedding
		text := PrepareTextForEmbedding(stage.Name, stage.Description)

		fmt.Printf("   [%d/%d] Generating embedding for: %s\n", i+1, totalStages, stage.Name)

		// Generate embedding
		embedding, err := e.client.GenerateEmbedding(text)
		if err != nil {
			fmt.Printf("   ⚠️  Warning: Failed for %s: %v\n", stage.ID, err)
			failedCount++
			continue
		}

		stage.Embedding = embedding
		updatedCount++

		// Small delay to avoid rate limits
		time.Sleep(100 * time.Millisecond)

		updatedStages[i] = stage
	}

	fmt.Printf("\n   ✅ Generated: %d | ⏭️  Skipped: %d | ❌ Failed: %d\n", updatedCount, skippedCount, failedCount)

	// Save updated data back to file
	if updatedCount > 0 {
		fmt.Println("   💾 Saving stages data...")
		e.stagesStorage.Set(&entities.StagesData{Stages: updatedStages})
		if err := e.stagesStorage.Save(); err != nil {
			return err
		}
		fmt.Println("   ✅ Stages data saved successfully")
	} else {
		fmt.Println("   ℹ️  No new embeddings to save")
	}

	return nil
}

// CosineSimilarity calculates the cosine similarity between two vectors
// Returns a value between -1 and 1, where 1 means identical, 0 means orthogonal, -1 means opposite
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// PrepareTextForEmbedding combines name and description for better search results
func PrepareTextForEmbedding(name, description string) string {
	return fmt.Sprintf("Name: %s\nDescription: %s", name, description)
}

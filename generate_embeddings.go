package main

import (
	"fmt"
	"time"
)

// GenerateAllEmbeddings generates embeddings for all models and stages
// This should be run when new content is added or descriptions change
func GenerateAllEmbeddings() error {
	fmt.Println("\n========================================")
	fmt.Println("🚀 Starting AI Embedding Generation")
	fmt.Println("========================================\n")

	// Generate embeddings for models
	fmt.Println("📦 Processing Models...")
	if err := GenerateModelsEmbeddings(); err != nil {
		return fmt.Errorf("error generating model embeddings: %w", err)
	}

	fmt.Println()

	// Generate embeddings for stages
	fmt.Println("🎭 Processing Stages...")
	if err := GenerateStagesEmbeddings(); err != nil {
		return fmt.Errorf("error generating stage embeddings: %w", err)
	}

	fmt.Println("\n========================================")
	fmt.Println("✅ Embedding Generation Complete!")
	fmt.Println("🔍 Search is now enabled")
	fmt.Println("========================================\n")
	return nil
}

// GenerateModelsEmbeddings generates embeddings for all models
func GenerateModelsEmbeddings() error {
	// Load existing models data
	modelsData, err := LoadModelsData()
	if err != nil {
		return err
	}

	totalModels := len(modelsData.Models)
	skippedCount := 0
	updatedCount := 0
	failedCount := 0

	fmt.Printf("   Found %d models total\n", totalModels)

	for i := range modelsData.Models {
		model := &modelsData.Models[i]

		// Skip if embedding already exists
		if len(model.Embedding) > 0 {
			skippedCount++
			continue
		}

		// Prepare text for embedding
		text := PrepareTextForEmbedding(model.Name, model.Description)

		fmt.Printf("   [%d/%d] Generating embedding for: %s\n", i+1, totalModels, model.Name)

		// Generate embedding
		embedding, err := GenerateEmbedding(text)
		if err != nil {
			fmt.Printf("   ⚠️  Warning: Failed for %s: %v\n", model.ID, err)
			failedCount++
			continue
		}

		model.Embedding = embedding
		updatedCount++

		// Small delay to avoid rate limits (optional, adjust based on your API tier)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\n   ✅ Generated: %d | ⏭️  Skipped: %d | ❌ Failed: %d\n", updatedCount, skippedCount, failedCount)

	// Save updated data back to file
	if updatedCount > 0 {
		fmt.Println("   💾 Saving models data...")
		if err := SaveModelsData(modelsData); err != nil {
			return err
		}
		fmt.Println("   ✅ Models data saved successfully")
	} else {
		fmt.Println("   ℹ️  No new embeddings to save")
	}

	return nil
}

// GenerateStagesEmbeddings generates embeddings for all stages
func GenerateStagesEmbeddings() error {
	// Load existing stages data
	stagesData, err := LoadStagesData()
	if err != nil {
		return err
	}

	totalStages := len(stagesData.Stages)
	skippedCount := 0
	updatedCount := 0
	failedCount := 0

	fmt.Printf("   Found %d stages total\n", totalStages)

	for i := range stagesData.Stages {
		stage := &stagesData.Stages[i]

		// Skip if embedding already exists
		if len(stage.Embedding) > 0 {
			skippedCount++
			continue
		}

		// Prepare text for embedding
		text := PrepareTextForEmbedding(stage.Name, stage.Description)

		fmt.Printf("   [%d/%d] Generating embedding for: %s\n", i+1, totalStages, stage.Name)

		// Generate embedding
		embedding, err := GenerateEmbedding(text)
		if err != nil {
			fmt.Printf("   ⚠️  Warning: Failed for %s: %v\n", stage.ID, err)
			failedCount++
			continue
		}

		stage.Embedding = embedding
		updatedCount++

		// Small delay to avoid rate limits
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\n   ✅ Generated: %d | ⏭️  Skipped: %d | ❌ Failed: %d\n", updatedCount, skippedCount, failedCount)

	// Save updated data back to file
	if updatedCount > 0 {
		fmt.Println("   💾 Saving stages data...")
		if err := SaveStagesData(stagesData); err != nil {
			return err
		}
		fmt.Println("   ✅ Stages data saved successfully")
	} else {
		fmt.Println("   ℹ️  No new embeddings to save")
	}

	return nil
}

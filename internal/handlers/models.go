package handlers

import (
	"fmt"
	"sort"

	"MMDContent/internal/entities"
	"MMDContent/internal/services/openai"
	"MMDContent/internal/storage"
)

type Models struct {
	client        openai.Client
	modelsStorage *storage.Models
}

func NewModels(
	client openai.Client,
	modelsStorage *storage.Models,
) *Models {
	return &Models{
		client:        client,
		modelsStorage: modelsStorage,
	}
}

// SearchModels searches models using semantic similarity with embeddings
func (a *Models) SearchModels(query string, limit int) ([]entities.Model, error) {
	if a.modelsStorage.IsEmpty() {
		return []entities.Model{}, nil
	}

	// Generate embedding for the search query
	queryEmbedding, err := a.client.GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Calculate similarity scores for all models
	type scoredModel struct {
		model entities.Model
		score float64
	}

	var scoredModels []scoredModel
	for _, model := range a.modelsStorage.Get().Models {
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

	results := make([]entities.Model, limit)
	for i := 0; i < limit; i++ {
		results[i] = scoredModels[i].model
	}

	return results, nil
}

// GetModels returns paginated models
func (a *Models) GetModels(page, perPage int) entities.PaginatedModels {
	if a.modelsStorage.IsEmpty() {
		return entities.PaginatedModels{
			Models:     []entities.Model{},
			Total:      0,
			Page:       page,
			PerPage:    perPage,
			TotalPages: 0,
		}
	}

	return a.modelsStorage.Get().GetPaginatedModels(page, perPage)
}

// GetAllModels returns all models without pagination
func (a *Models) GetAllModels() []entities.Model {
	if a.modelsStorage.IsEmpty() {
		return []entities.Model{}
	}

	return a.modelsStorage.Get().Models
}

// RefreshModelsData re-parses and reloads the models data
func (a *Models) RefreshModelsData() error {
	err := entities.ParseModelsData()
	if err != nil {
		return err
	}

	modelsData, err := entities.LoadModelsData()
	if err != nil {
		return err
	}

	a.modelsStorage.Set(modelsData)
	return nil
}

package handlers

import (
	"fmt"
	"sort"

	"MMDContent/internal/entities"
	"MMDContent/internal/services/openai"
	"MMDContent/internal/storage"
)

type Stages struct {
	client        openai.Client
	stagesStorage *storage.Stages
}

func NewStages(
	client openai.Client,
	stagesStorage *storage.Stages,
) *Stages {
	return &Stages{
		client:        client,
		stagesStorage: stagesStorage,
	}
}

// SearchStages searches stages using semantic similarity with embeddings
func (a *Stages) SearchStages(query string, limit int) ([]entities.Stage, error) {
	if a.stagesStorage.IsEmpty() {
		return []entities.Stage{}, nil
	}

	// Generate embedding for the search query
	queryEmbedding, err := a.client.GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Calculate similarity scores for all stages
	type scoredStage struct {
		stage entities.Stage
		score float64
	}

	var scoredStages []scoredStage
	for _, stage := range a.stagesStorage.Get().Stages {
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

	results := make([]entities.Stage, limit)
	for i := 0; i < limit; i++ {
		results[i] = scoredStages[i].stage
	}

	return results, nil
}

// GetAllStages returns all stages without pagination
func (a *Stages) GetAllStages() []entities.Stage {
	if a.stagesStorage.IsEmpty() {
		return []entities.Stage{}
	}

	return a.stagesStorage.Get().Stages
}

// GetStages returns paginated stages
func (a *Stages) GetStages(page, perPage int) entities.Pagination[entities.Stage] {
	if a.stagesStorage.IsEmpty() {
		return entities.Pagination[entities.Stage]{
			Data:       []entities.Stage{},
			Total:      0,
			Page:       page,
			PerPage:    perPage,
			TotalPages: 0,
		}
	}

	return a.stagesStorage.GetPaginatedStages(page, perPage)
}

// RefreshStagesData re-parses and reloads the stages data
func (a *Stages) RefreshStagesData() error {
	err := a.stagesStorage.Refresh()
	if err != nil {
		return err
	}

	return nil
}

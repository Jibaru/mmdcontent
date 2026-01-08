package handlers

import (
	"fmt"
	"sort"

	"MMDContent/internal/entities"
	"MMDContent/internal/services/openai"
	"MMDContent/internal/storage"
)

type Motions struct {
	client         openai.Client
	motionsStorage *storage.Motions
}

func NewMotions(
	client openai.Client,
	motionsStorage *storage.Motions,
) *Motions {
	return &Motions{
		client:         client,
		motionsStorage: motionsStorage,
	}
}

// SearchMotions searches motions using semantic similarity with embeddings
func (a *Motions) SearchMotions(query string, limit int) ([]entities.Motion, error) {
	if a.motionsStorage.IsEmpty() {
		return []entities.Motion{}, nil
	}

	// Generate embedding for the search query
	queryEmbedding, err := a.client.GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Calculate similarity scores for all motions
	type scoredMotion struct {
		motion entities.Motion
		score  float64
	}

	var scoredMotions []scoredMotion
	for _, motion := range a.motionsStorage.Get().Motions {
		if len(motion.Embedding) == 0 {
			// Skip motions without embeddings
			continue
		}

		similarity := CosineSimilarity(queryEmbedding, motion.Embedding)
		scoredMotions = append(scoredMotions, scoredMotion{
			motion: motion,
			score:  similarity,
		})
	}

	// Sort by similarity score (descending)
	sort.Slice(scoredMotions, func(i, j int) bool {
		return scoredMotions[i].score > scoredMotions[j].score
	})

	// Return top results
	if limit <= 0 || limit > len(scoredMotions) {
		limit = len(scoredMotions)
	}

	results := make([]entities.Motion, limit)
	for i := 0; i < limit; i++ {
		results[i] = scoredMotions[i].motion
	}

	return results, nil
}

// GetMotions returns paginated motions
func (a *Motions) GetMotions(page, perPage int) entities.Pagination[entities.Motion] {
	if a.motionsStorage.IsEmpty() {
		return entities.Pagination[entities.Motion]{
			Data:       []entities.Motion{},
			Total:      0,
			Page:       page,
			PerPage:    perPage,
			TotalPages: 0,
		}
	}

	return a.motionsStorage.GetPaginatedMotions(page, perPage)
}

// GetAllMotions returns all motions without pagination
func (a *Motions) GetAllMotions() []entities.Motion {
	if a.motionsStorage.IsEmpty() {
		return []entities.Motion{}
	}

	return a.motionsStorage.Get().Motions
}

// RefreshMotionsData re-parses and reloads the motions data
func (a *Motions) RefreshMotionsData() error {
	err := a.motionsStorage.Refresh()
	if err != nil {
		return err
	}

	return nil
}

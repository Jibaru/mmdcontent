package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Model struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Screenshots  []string `json:"screenshots"`
	Description  string   `json:"description"`
	OriginalPath string   `json:"originalPath"`
}

type ModelsData struct {
	Models []Model `json:"models"`
}

type PaginatedModels struct {
	Models      []Model `json:"models"`
	Total       int     `json:"total"`
	Page        int     `json:"page"`
	PerPage     int     `json:"perPage"`
	TotalPages  int     `json:"totalPages"`
}

// ParseModelsData reads the data/Models directory and creates a data.json file
func ParseModelsData() error {
	modelsDir := "data/Models"
	var models []Model

	// Read all directories in data/Models
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modelID := entry.Name()
		modelPath := filepath.Join(modelsDir, modelID)

		// Read ruta.txt
		rutaPath := filepath.Join(modelPath, "ruta.txt")
		rutaContent, err := os.ReadFile(rutaPath)
		if err != nil {
			continue // Skip if ruta.txt doesn't exist
		}

		// Extract filename from path
		rutaStr := strings.TrimSpace(string(rutaContent))
		modelName := filepath.Base(rutaStr)

		// Read descripcion.txt
		descPath := filepath.Join(modelPath, "descripcion.txt")
		descContent, err := os.ReadFile(descPath)
		description := ""
		if err == nil {
			description = string(descContent)
		}

		// Read screenshots
		screenshotsDir := filepath.Join(modelPath, "screenshots")
		var screenshots []string
		screenshotEntries, err := os.ReadDir(screenshotsDir)
		if err == nil {
			for _, se := range screenshotEntries {
				if !se.IsDir() {
					// Store absolute path
					absPath, err := filepath.Abs(filepath.Join(screenshotsDir, se.Name()))
					if err == nil {
						screenshots = append(screenshots, absPath)
					}
				}
			}
		}

		// Sort screenshots
		sort.Strings(screenshots)

		model := Model{
			ID:           modelID,
			Name:         modelName,
			Screenshots:  screenshots,
			Description:  description,
			OriginalPath: rutaStr,
		}

		models = append(models, model)
	}

	// Sort models by ID
	sort.Slice(models, func(i, j int) bool {
		return models[i].ID < models[j].ID
	})

	// Create ModelsData struct
	modelsData := ModelsData{
		Models: models,
	}

	// Write to data/data.json
	jsonData, err := json.MarshalIndent(modelsData, "", "  ")
	if err != nil {
		return err
	}

	dataPath := filepath.Join("data", "data.json")
	err = os.WriteFile(dataPath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// LoadModelsData loads the models from data/data.json
func LoadModelsData() (*ModelsData, error) {
	dataPath := filepath.Join("data", "data.json")
	jsonData, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var modelsData ModelsData
	err = json.Unmarshal(jsonData, &modelsData)
	if err != nil {
		return nil, err
	}

	return &modelsData, nil
}

// GetPaginatedModels returns a paginated subset of models
func (md *ModelsData) GetPaginatedModels(page, perPage int) PaginatedModels {
	total := len(md.Models)

	// Calculate pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 100
	}

	totalPages := (total + perPage - 1) / perPage
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * perPage
	end := start + perPage

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedModels := md.Models[start:end]

	return PaginatedModels{
		Models:     paginatedModels,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}
}

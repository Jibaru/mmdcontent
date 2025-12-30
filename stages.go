package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Stage struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Screenshots  []string `json:"screenshots"`
	Description  string   `json:"description"`
	OriginalPath string   `json:"originalPath"`
}

type StagesData struct {
	Stages []Stage `json:"stages"`
}

type PaginatedStages struct {
	Stages      []Stage `json:"stages"`
	Total       int     `json:"total"`
	Page        int     `json:"page"`
	PerPage     int     `json:"perPage"`
	TotalPages  int     `json:"totalPages"`
}

// ParseStagesData reads the data/Stages directory and creates a data/stages.json file
func ParseStagesData() error {
	stagesDir := "data/Stages"
	var stages []Stage

	// Read all directories in data/Stages
	entries, err := os.ReadDir(stagesDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		stageID := entry.Name()
		stagePath := filepath.Join(stagesDir, stageID)

		// Read ruta.txt
		rutaPath := filepath.Join(stagePath, "ruta.txt")
		rutaContent, err := os.ReadFile(rutaPath)
		if err != nil {
			continue // Skip if ruta.txt doesn't exist
		}

		// Extract filename from path
		rutaStr := strings.TrimSpace(string(rutaContent))
		stageName := filepath.Base(rutaStr)

		// Read descripcion.txt
		descPath := filepath.Join(stagePath, "descripcion.txt")
		descContent, err := os.ReadFile(descPath)
		description := ""
		if err == nil {
			description = string(descContent)
		}

		// Read screenshots
		screenshotsDir := filepath.Join(stagePath, "screenshots")
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

		stage := Stage{
			ID:           stageID,
			Name:         stageName,
			Screenshots:  screenshots,
			Description:  description,
			OriginalPath: rutaStr,
		}

		stages = append(stages, stage)
	}

	// Sort stages by ID
	sort.Slice(stages, func(i, j int) bool {
		return stages[i].ID < stages[j].ID
	})

	// Create StagesData struct
	stagesData := StagesData{
		Stages: stages,
	}

	// Write to data/stages.json
	jsonData, err := json.MarshalIndent(stagesData, "", "  ")
	if err != nil {
		return err
	}

	dataPath := filepath.Join("data", "stages.json")
	err = os.WriteFile(dataPath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// LoadStagesData loads the stages from data/stages.json
func LoadStagesData() (*StagesData, error) {
	dataPath := filepath.Join("data", "stages.json")
	jsonData, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var stagesData StagesData
	err = json.Unmarshal(jsonData, &stagesData)
	if err != nil {
		return nil, err
	}

	return &stagesData, nil
}

// GetPaginatedStages returns a paginated subset of stages
func (sd *StagesData) GetPaginatedStages(page, perPage int) PaginatedStages {
	total := len(sd.Stages)

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

	paginatedStages := sd.Stages[start:end]

	return PaginatedStages{
		Stages:     paginatedStages,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}
}

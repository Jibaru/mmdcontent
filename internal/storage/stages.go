package storage

import (
	"MMDContent/internal/entities"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Stages struct {
	data *entities.StagesData
}

func NewStages() *Stages {
	return &Stages{
		data: &entities.StagesData{},
	}
}

func (m *Stages) Get() *entities.StagesData {
	return m.data
}

func (m *Stages) Set(data *entities.StagesData) {
	m.data = data
}

func (m *Stages) IsEmpty() bool {
	return len(m.data.Stages) == 0
}

func (m *Stages) Total() int {
	return len(m.data.Stages)
}

func (m *Stages) Refresh() error {
	err := ParseStagesData()
	if err != nil {
		return err
	}

	data, err := LoadStagesData()
	if err != nil {
		return err
	}

	m.data = data
	return nil
}

func (m *Stages) Save() error {
	return SaveStagesData(m.data)
}

// GetPaginatedStages returns a paginated subset of stages
func (m *Stages) GetPaginatedStages(page, perPage int) entities.Pagination[entities.Stage] {
	total := m.Total()

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

	paginatedStages := m.data.Stages[start:end]

	return entities.Pagination[entities.Stage]{
		Data:       paginatedStages,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}
}

// ParseStagesData reads the data/Stages directory and creates a data/stages.json file
func ParseStagesData() error {
	stagesDir := "data/Stages"
	var stages []entities.Stage

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

		stage := entities.Stage{
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
	stagesData := entities.StagesData{
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
func LoadStagesData() (*entities.StagesData, error) {
	dataPath := filepath.Join("data", "stages.json")
	jsonData, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var stagesData entities.StagesData
	err = json.Unmarshal(jsonData, &stagesData)
	if err != nil {
		return nil, err
	}

	return &stagesData, nil
}

// SaveStagesData saves the stages data to data/stages.json
func SaveStagesData(stagesData *entities.StagesData) error {
	dataPath := filepath.Join("data", "stages.json")
	jsonData, err := json.MarshalIndent(stagesData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(dataPath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

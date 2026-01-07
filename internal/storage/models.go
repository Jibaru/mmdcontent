package storage

import (
	"MMDContent/internal/entities"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Models struct {
	data *entities.ModelsData
}

func NewModels() *Models {
	return &Models{
		data: &entities.ModelsData{},
	}
}

func (m *Models) Get() *entities.ModelsData {
	return m.data
}

func (m *Models) Set(data *entities.ModelsData) {
	m.data = data
}

func (m *Models) IsEmpty() bool {
	return len(m.data.Models) == 0
}

func (m *Models) Total() int {
	return len(m.data.Models)
}

func (m *Models) Refresh() error {
	err := ParseModelsData()
	if err != nil {
		return err
	}

	modelsData, err := LoadModelsData()
	if err != nil {
		return err
	}

	m.data = modelsData
	return nil
}

func (m *Models) Save() error {
	return SaveModelsData(m.data)
}

// GetPaginatedModels returns a paginated subset of models
func (m *Models) GetPaginatedModels(page, perPage int) entities.Pagination[entities.Model] {
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

	paginatedModels := m.data.Models[start:end]

	return entities.Pagination[entities.Model]{
		Data:       paginatedModels,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}
}

// ParseModelsData reads the data/Models directory and creates a data.json file
func ParseModelsData() error {
	modelsDir := "data/Models"
	var models []entities.Model

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

		model := entities.Model{
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
	modelsData := entities.ModelsData{
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
func LoadModelsData() (*entities.ModelsData, error) {
	dataPath := filepath.Join("data", "data.json")
	jsonData, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var modelsData entities.ModelsData
	err = json.Unmarshal(jsonData, &modelsData)
	if err != nil {
		return nil, err
	}

	return &modelsData, nil
}

// SaveModelsData saves the models data to data/data.json
func SaveModelsData(modelsData *entities.ModelsData) error {
	dataPath := filepath.Join("data", "data.json")
	jsonData, err := json.MarshalIndent(modelsData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(dataPath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

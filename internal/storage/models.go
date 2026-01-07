package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"MMDContent/internal/entities"
)

type Models struct {
	data     *entities.ModelsData
	dirName  string
	filename string
}

func NewModelsLoaded(dirName string, filename string) (*Models, error) {
	m := &Models{
		dirName:  dirName,
		filename: filename,
	}

	err := m.sync()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Models) sync() error {
	modelsDataInFolder, err := readModelsDataFromFolder(m.dirName)
	if err != nil {
		return err
	}

	modelsDataInJSON, err := loadModelsDataFromFile(m.filename)
	if err != nil {
		return err
	}

	missingModels := make([]entities.Model, 0)
	for _, model := range modelsDataInFolder.Models {
		if !modelsDataInJSON.Has(model) {
			missingModels = append(missingModels, model)
		}
	}

	modelsDataInJSON.Models = append(modelsDataInJSON.Models, missingModels...)
	m.Set(modelsDataInJSON)
	return m.Save()
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
	return m.sync()
}

func (m *Models) Save() error {
	jsonData, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(m.filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
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

func loadModelsDataFromFile(filename string) (*entities.ModelsData, error) {
	jsonData, err := os.ReadFile(filename)
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

func readModelsDataFromFolder(dirName string) (*entities.ModelsData, error) {
	var models []entities.Model

	// Read all directories in data/Models
	entries, err := os.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modelID := entry.Name()
		modelPath := filepath.Join(dirName, modelID)

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
	return &entities.ModelsData{
		Models: models,
	}, nil
}

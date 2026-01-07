package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"MMDContent/internal/entities"
)

type Stages struct {
	data     *entities.StagesData
	dirName  string
	filename string
}

func NewStagesLoaded(dirName string, filename string) (*Stages, error) {
	s := &Stages{
		dirName:  dirName,
		filename: filename,
	}

	err := s.sync()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *Stages) sync() error {
	stagesDataInFolder, err := readStagesDataFromFolder(m.dirName)
	if err != nil {
		return err
	}

	stagesDataInJSON, err := loadStagesDataFromFile(m.filename)
	if err != nil {
		return err
	}

	missingStages := make([]entities.Stage, 0)
	for _, stage := range stagesDataInFolder.Stages {
		if !stagesDataInJSON.Has(stage) {
			missingStages = append(missingStages, stage)
		}
	}

	stagesDataInJSON.Stages = append(stagesDataInJSON.Stages, missingStages...)
	m.Set(stagesDataInJSON)
	return m.Save()
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
	return m.sync()
}

func (m *Stages) Save() error {
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

func readStagesDataFromFolder(dirName string) (*entities.StagesData, error) {
	var stages []entities.Stage

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

		stage := entities.Stage{
			ID:           modelID,
			Name:         modelName,
			Screenshots:  screenshots,
			Description:  description,
			OriginalPath: rutaStr,
		}

		stages = append(stages, stage)
	}

	// Sort models by ID
	sort.Slice(stages, func(i, j int) bool {
		return stages[i].ID < stages[j].ID
	})

	// Create ModelsData struct
	return &entities.StagesData{
		Stages: stages,
	}, nil
}

func loadStagesDataFromFile(filename string) (*entities.StagesData, error) {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data entities.StagesData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

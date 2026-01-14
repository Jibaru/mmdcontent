package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"MMDContent/internal/entities"
)

type Motions struct {
	data     *entities.MotionsData
	dirName  string
	filename string
}

func NewMotionsLoaded(dirName string, filename string) (*Motions, error) {
	m := &Motions{
		dirName:  dirName,
		filename: filename,
	}

	err := m.sync()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Motions) sync() error {
	motionsDataInFolder, err := readMotionsDataFromFolder(m.dirName)
	if err != nil {
		return err
	}

	motionsDataInJSON, err := loadMotionsDataFromFile(m.filename)
	if err != nil {
		return err
	}

	missingMotions := make([]entities.Motion, 0)
	for _, motion := range motionsDataInFolder.Motions {
		if !motionsDataInJSON.Has(motion) {
			missingMotions = append(missingMotions, motion)
		}
	}

	motionsDataInJSON.Motions = append(motionsDataInJSON.Motions, missingMotions...)
	m.Set(motionsDataInJSON)
	return m.Save()
}

func (m *Motions) Get() *entities.MotionsData {
	return m.data
}

func (m *Motions) Set(data *entities.MotionsData) {
	m.data = data
}

func (m *Motions) IsEmpty() bool {
	return len(m.data.Motions) == 0
}

func (m *Motions) Total() int {
	return len(m.data.Motions)
}

func (m *Motions) Refresh() error {
	return m.sync()
}

func (m *Motions) Save() error {
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

// GetPaginatedMotions returns a paginated subset of motions
func (m *Motions) GetPaginatedMotions(page, perPage int) entities.Pagination[entities.Motion] {
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

	paginatedMotions := m.data.Motions[start:end]

	return entities.Pagination[entities.Motion]{
		Data:       paginatedMotions,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}
}

func loadMotionsDataFromFile(filename string) (*entities.MotionsData, error) {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var motionsData entities.MotionsData
	err = json.Unmarshal(jsonData, &motionsData)
	if err != nil {
		return nil, err
	}

	return &motionsData, nil
}

func readMotionsDataFromFolder(dirName string) (*entities.MotionsData, error) {
	var motions []entities.Motion

	// Read all directories in data/Motions
	entries, err := os.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		motionID := entry.Name()
		motionPath := filepath.Join(dirName, motionID)

		// Read ruta.txt
		rutaPath := filepath.Join(motionPath, "ruta.txt")
		rutaContent, err := os.ReadFile(rutaPath)
		if err != nil {
			continue // Skip if ruta.txt doesn't exist
		}

		// Extract filename from path
		rutaStr := strings.TrimSpace(string(rutaContent))
		motionName := filepath.Base(rutaStr)

		// Read descripcion.txt
		descPath := filepath.Join(motionPath, "descripcion.txt")
		descContent, err := os.ReadFile(descPath)
		description := ""
		if err == nil {
			description = string(descContent)
		}

		// Read screenshots
		screenshotsDir := filepath.Join(motionPath, "screenshots")
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

		// Read video files
		videoDir := filepath.Join(motionPath, "video")
		var video []string
		videoEntries, err := os.ReadDir(videoDir)
		if err == nil {
			for _, ve := range videoEntries {
				if !ve.IsDir() {
					// Store absolute path
					absPath, err := filepath.Abs(filepath.Join(videoDir, ve.Name()))
					if err == nil {
						video = append(video, absPath)
					}
				}
			}
		}

		// Sort video files
		sort.Strings(video)

		motion := entities.Motion{
			ID:           motionID,
			Name:         motionName,
			Screenshots:  screenshots,
			Video:        video,
			Description:  description,
			OriginalPath: rutaStr,
		}

		motions = append(motions, motion)
	}

	// Sort motions by ID
	sort.Slice(motions, func(i, j int) bool {
		return motions[i].ID < motions[j].ID
	})

	// Create MotionsData struct
	return &entities.MotionsData{
		Motions: motions,
	}, nil
}

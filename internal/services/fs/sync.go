package fs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"MMDContent/internal/entities"
)

// ChangeType represents the type of change detected
type ChangeType int

const (
	ChangeAdded ChangeType = iota
	ChangeModified
	ChangeDeleted
)

// Change represents a detected change in the file system
type Change struct {
	ID   string
	Type ChangeType
	Path string
}

// CheckModelsChanges detects changes in the data/Models directory
func CheckModelsChanges() ([]Change, error) {
	return checkDirectoryChanges("data/Models", "data/data.json")
}

// CheckStagesChanges detects changes in the data/Stages directory
func CheckStagesChanges() ([]Change, error) {
	return checkDirectoryChanges("data/Stages", "data/stages.json")
}

// checkDirectoryChanges compares directory contents with JSON file
func checkDirectoryChanges(dirPath, jsonPath string) ([]Change, error) {
	var changes []Change

	// Get current directory state
	currentDirs, err := getDirectoryIDs(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	// Get JSON file state
	jsonDirs, err := getJSONIDs(jsonPath)
	if err != nil {
		// If JSON doesn't exist, all current dirs are new
		if os.IsNotExist(err) {
			for id := range currentDirs {
				changes = append(changes, Change{
					ID:   id,
					Type: ChangeAdded,
					Path: filepath.Join(dirPath, id),
				})
			}
			return changes, nil
		}
		return nil, fmt.Errorf("error reading JSON: %w", err)
	}

	// Check for new or modified directories
	for id, currentModTime := range currentDirs {
		jsonModTime, existsInJSON := jsonDirs[id]

		if !existsInJSON {
			// New directory
			changes = append(changes, Change{
				ID:   id,
				Type: ChangeAdded,
				Path: filepath.Join(dirPath, id),
			})
		} else if currentModTime.After(jsonModTime) {
			// Modified directory (content changed)
			changes = append(changes, Change{
				ID:   id,
				Type: ChangeModified,
				Path: filepath.Join(dirPath, id),
			})
		}
	}

	// Check for deleted directories
	for id := range jsonDirs {
		if _, existsInDir := currentDirs[id]; !existsInDir {
			changes = append(changes, Change{
				ID:   id,
				Type: ChangeDeleted,
				Path: filepath.Join(dirPath, id),
			})
		}
	}

	return changes, nil
}

// getDirectoryIDs returns a map of directory IDs to their modification times
func getDirectoryIDs(dirPath string) (map[string]time.Time, error) {
	dirs := make(map[string]time.Time)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Get the most recent modification time from the directory contents
		dirFullPath := filepath.Join(dirPath, entry.Name())
		modTime, err := getDirectoryModTime(dirFullPath)
		if err != nil {
			fmt.Printf("Warning: Could not get mod time for %s: %v\n", entry.Name(), err)
			continue
		}

		dirs[entry.Name()] = modTime
	}

	return dirs, nil
}

// getDirectoryModTime returns the most recent modification time from a directory and its contents
func getDirectoryModTime(dirPath string) (time.Time, error) {
	var latestModTime time.Time

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.ModTime().After(latestModTime) {
			latestModTime = info.ModTime()
		}

		return nil
	})

	if err != nil {
		return time.Time{}, err
	}

	return latestModTime, nil
}

// getJSONIDs returns a map of IDs from the JSON file with their last known mod times
// This is a simplified version that checks if the item exists in JSON
func getJSONIDs(jsonPath string) (map[string]time.Time, error) {
	ids := make(map[string]time.Time)

	// Read JSON file
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	// Get JSON file modification time as a baseline
	fileInfo, err := os.Stat(jsonPath)
	if err != nil {
		return nil, err
	}
	baseModTime := fileInfo.ModTime()

	// Parse JSON based on file type
	if filepath.Base(jsonPath) == "data.json" {
		var modelsData entities.ModelsData
		if err := json.Unmarshal(jsonData, &modelsData); err != nil {
			return nil, err
		}

		for _, model := range modelsData.Models {
			// Use JSON file mod time as proxy for when this was last synced
			ids[model.ID] = baseModTime
		}
	} else if filepath.Base(jsonPath) == "stages.json" {
		var stagesData entities.StagesData
		if err := json.Unmarshal(jsonData, &stagesData); err != nil {
			return nil, err
		}

		for _, stage := range stagesData.Stages {
			ids[stage.ID] = baseModTime
		}
	}

	return ids, nil
}

// NeedsSync checks if any changes are detected
func NeedsSync(changes []Change) bool {
	return len(changes) > 0
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

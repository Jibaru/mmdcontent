package fs

import (
	"fmt"

	"MMDContent/internal/entities"
	"MMDContent/internal/storage"
)

func SyncFiles() (*entities.ModelsData, *entities.StagesData, error) {
	// Check for models changes
	modelsChanges, err := CheckModelsChanges()
	if err != nil {
		return nil, nil, err
	}

	// Parse models data if changes detected or JSON doesn't exist
	if NeedsSync(modelsChanges) || !fileExists("data/data.json") {
		err := storage.ParseModelsData()
		if err != nil {
			return nil, nil, err
		}
	}

	// Load models data
	modelsData, err := storage.LoadModelsData()
	if err != nil {
		return nil, nil, err
	}

	// Check for stages changes
	fmt.Println("🎭 Checking Stages directory...")
	stagesChanges, err := CheckStagesChanges()
	if err != nil {
		return nil, nil, err
	}

	// Parse stages data if changes detected or JSON doesn't exist
	if NeedsSync(stagesChanges) || !fileExists("data/stages.json") {
		err = storage.ParseStagesData()
		if err != nil {
			return nil, nil, err
		}
	}

	// Load stages data
	stagesData, err := storage.LoadStagesData()
	if err != nil {
		return nil, nil, err
	}

	return modelsData, stagesData, nil
}

package storage

import "MMDContent/internal/entities"

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

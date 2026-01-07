package storage

import "MMDContent/internal/entities"

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

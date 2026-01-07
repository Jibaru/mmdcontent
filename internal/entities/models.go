package entities

type Model struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Screenshots  []string  `json:"screenshots"`
	Description  string    `json:"description"`
	OriginalPath string    `json:"originalPath"`
	Embedding    []float64 `json:"embedding,omitempty"`
}

func (m *Model) Equal(o Model) bool {
	return m.ID == o.ID && m.Name == o.Name &&
		m.Description == o.Description &&
		m.OriginalPath == o.OriginalPath &&
		len(m.Screenshots) == len(o.Screenshots) // TODO: improve here
}

type ModelsData struct {
	Models []Model `json:"models"`
}

func (m *ModelsData) Has(o Model) bool {
	for _, model := range m.Models {
		if model.Equal(o) {
			return true
		}
	}
	return false
}

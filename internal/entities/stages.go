package entities

type Stage struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Screenshots  []string  `json:"screenshots"`
	Description  string    `json:"description"`
	OriginalPath string    `json:"originalPath"`
	Embedding    []float64 `json:"embedding,omitempty"`
}

func (m *Stage) Equal(o Stage) bool {
	return m.ID == o.ID && m.Name == o.Name &&
		m.Description == o.Description &&
		m.OriginalPath == o.OriginalPath &&
		len(m.Screenshots) == len(o.Screenshots) // TODO: improve here
}

type StagesData struct {
	Stages []Stage `json:"stages"`
}

func (m *StagesData) Has(o Stage) bool {
	for _, stage := range m.Stages {
		if stage.Equal(o) {
			return true
		}
	}
	return false
}

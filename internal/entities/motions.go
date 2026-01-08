package entities

type Motion struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Screenshots  []string  `json:"screenshots"`
	Description  string    `json:"description"`
	OriginalPath string    `json:"originalPath"`
	Embedding    []float64 `json:"embedding,omitempty"`
}

func (m *Motion) Equal(o Motion) bool {
	return m.ID == o.ID && m.Name == o.Name &&
		m.Description == o.Description &&
		m.OriginalPath == o.OriginalPath &&
		equalSlices(m.Screenshots, o.Screenshots)
}

type MotionsData struct {
	Motions []Motion `json:"motions"`
}

func (m *MotionsData) Has(o Motion) bool {
	for _, motion := range m.Motions {
		if motion.Equal(o) {
			return true
		}
	}
	return false
}

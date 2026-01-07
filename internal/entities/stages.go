package entities

type Stage struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Screenshots  []string  `json:"screenshots"`
	Description  string    `json:"description"`
	OriginalPath string    `json:"originalPath"`
	Embedding    []float64 `json:"embedding,omitempty"`
}

type StagesData struct {
	Stages []Stage `json:"stages"`
}

type PaginatedStages struct {
	Stages     []Stage `json:"stages"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	PerPage    int     `json:"perPage"`
	TotalPages int     `json:"totalPages"`
}

package entities

type Model struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Screenshots  []string  `json:"screenshots"`
	Description  string    `json:"description"`
	OriginalPath string    `json:"originalPath"`
	Embedding    []float64 `json:"embedding,omitempty"`
}

type ModelsData struct {
	Models []Model `json:"models"`
}

type PaginatedModels struct {
	Models     []Model `json:"models"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	PerPage    int     `json:"perPage"`
	TotalPages int     `json:"totalPages"`
}

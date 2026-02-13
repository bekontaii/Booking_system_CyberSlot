package club

type Club struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	City     string `json:"city"`
	Address  string `json:"address"`
	IsActive bool   `json:"is_active"`
}

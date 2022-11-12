package inveli

type LimitPaylater struct {
	MaxLimit       float64 `json:"max_limit"`
	AvailableLimit float64 `json:"available_limit"`
}

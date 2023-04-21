package response

type GetRiwayatPaylaterPerbulanResponse struct {
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	Month      int     `json:"month"`
	Status     string  `json:"status_pembayaran"`
	TotalBayar float64 `json:"total_belanja"`
}

type GetListRiwayatPaylaterPerbulanResponse struct {
	Month      string  `json:"month"`
	Status     string  `json:"status_pembayaran"`
	TotalBayar float64 `json:"total_belanja"`
}

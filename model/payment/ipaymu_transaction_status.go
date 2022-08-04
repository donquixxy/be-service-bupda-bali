package payment

type PaymentStatusResponse struct {
	Status  int           `json:"Status"`
	Data    PaymentStatus `json:"Data"`
	Message string        `json:"Message"`
}

type PaymentStatus struct {
	TransactionId  int     `json:"TransactionId"`
	SessionId      string  `json:"SessionId"`
	ReferenceId    string  `json:"ReferenceId"`
	RelatedId      int     `json:"RelatedId"`
	Sender         string  `json:"Sender"`
	Recevier       string  `json:"Recevier"`
	Amount         float64 `json:"Amount"`
	Fee            float64 `json:"Fee"`
	Status         int     `json:"Status"`
	StatusDesc     string  `json:"StatusDesc"`
	Type           int     `json:"Type"`
	TypeDesc       string  `json:"TypeDesc"`
	Notes          string  `json:"Notes"`
	CreatedDate    string  `json:"CreatedDate"`
	ExpiredDate    string  `json:"ExpiredDate"`
	SuccessDate    string  `json:"SuccessDate"`
	SettlementDate string  `json:"SettlementDate"`
}

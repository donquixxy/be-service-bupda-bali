package payment

type PaymentStatusResponse struct {
	Status  int
	Data    PaymentStatus
	Message string
}

type PaymentStatus struct {
	TransactionId  int
	SessionId      string
	ReferenceId    string
	RelatedId      int
	Sender         string
	Recevier       string
	Amount         string
	Fee            string
	Status         int
	StatusDesc     string
	Type           int
	TypeDesc       string
	Notes          string
	CreatedDate    string
	ExpiredDate    string
	SuccessDate    string
	SettlementDate string
}

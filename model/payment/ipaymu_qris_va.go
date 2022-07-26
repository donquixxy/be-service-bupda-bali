package payment

type IpaymuQrisVaRequest struct {
	Name           string
	Phone          string
	Email          string
	Amount         float64
	ReferenceId    string
	PaymentMethod  string
	PaymentChannel string
}

type IpaymuQrisVaResponse struct {
	Status  int
	Message string
	Data    Data
}

type Data struct {
	SessionId     string
	TransactionId int
	ReferenceId   string
	Via           string
	Channel       string
	PaymentNo     string
	PaymentName   string
	Total         float64
	Fee           float64
	Expired       string
}

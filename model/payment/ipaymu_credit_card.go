package payment

type IpaymuCreditCardRequest struct {
	Product       []string
	Qty           []int
	Price         []float64
	ReturnUrl     string
	CancelUrl     string
	NotifyUrl     string
	ReferenceId   string
	BuyerName     string
	BuyerEmail    string
	BuyerPhone    string
	PaymentMethod string
}

type IpaymuCreditCardResponse struct {
	Status  int
	Message string
	Data    CCData
}

type CCData struct {
	SessionId string
	Url       string
}

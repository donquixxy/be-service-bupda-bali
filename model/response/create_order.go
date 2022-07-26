package response

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type CreateOrderResponse struct {
	Id             string  `json:"id_order"`
	NumberOrder    string  `json:"number_order"`
	PaymentNo      string  `json:"payment_no"`
	PaymentName    string  `json:"payment_name"`
	Total          float64 `json:"total"`
	Expired        string  `json:"expired"`
	PaymentMethod  string  `json:"payment_method"`
	PaymentChannel string  `json:"payment_channel"`
	PaymentStatus  int     `json:"payment_status"`
	BankName       string  `json:"bank_name"`
	BankLogo       string  `json:"bank_logo"`
}

func ToCreateOrderResponse(order *entity.Order, paymentChannel *entity.PaymentChannel) (orderResponse CreateOrderResponse) {
	orderResponse.Id = order.Id
	orderResponse.NumberOrder = order.NumberOrder
	orderResponse.PaymentNo = order.PaymentNo
	orderResponse.PaymentName = order.PaymentName
	orderResponse.Total = order.PaymentCash
	orderResponse.Expired = order.PaymentDueDate.Time.Format("2006-01-02 15:04:05")
	orderResponse.PaymentMethod = order.PaymentMethod
	orderResponse.PaymentChannel = order.PaymentChannel
	orderResponse.PaymentStatus = order.PaymentStatus
	orderResponse.BankName = paymentChannel.Name
	orderResponse.BankLogo = paymentChannel.Logo
	return orderResponse
}

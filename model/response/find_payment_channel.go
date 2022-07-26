package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindPaymentChannelResponse struct {
	Id                 string  `json:"id"`
	IdPaymentMethod    string  `json:"id_payment_method"`
	Name               string  `json:"name"`
	Code               string  `json:"payment_code"`
	MethodCode         string  `json:"payment_method"`
	Logo               string  `json:"logo"`
	AdminFee           float64 `json:"admin_fee"`
	AdminFeePercentage float64 `json:"admin_fee_percentage"`
}

func ToFindPaymentChannelResponse(paymentChannels []entity.PaymentChannel) (paymentChannelResponses []FindPaymentChannelResponse) {
	for _, paymentChannel := range paymentChannels {
		paymentChannelResponse := FindPaymentChannelResponse{}
		paymentChannelResponse.Id = paymentChannel.Id
		paymentChannelResponse.IdPaymentMethod = paymentChannel.IdPaymentMethod
		paymentChannelResponse.Name = paymentChannel.Name
		paymentChannelResponse.Code = paymentChannel.Code
		paymentChannelResponse.Logo = paymentChannel.Logo
		paymentChannelResponse.MethodCode = paymentChannel.PaymentMethod.MethodCode
		paymentChannelResponse.AdminFee = paymentChannel.AdminFee
		paymentChannelResponse.AdminFeePercentage = paymentChannel.AdminFeePercentage
		paymentChannelResponses = append(paymentChannelResponses, paymentChannelResponse)
	}
	return paymentChannelResponses
}

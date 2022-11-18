package response

import (
	"log"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

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

func ToFindPaymentChannelResponse(paymentChannels []entity.PaymentChannel, statusUser int, biayaTanggungRenteng float64, isPaylater int) (paymentChannelResponses []FindPaymentChannelResponse) {
	for _, paymentChannel := range paymentChannels {
		paymentChannelResponse := FindPaymentChannelResponse{}
		log.Println("payment channel", paymentChannel.Code)
		if paymentChannel.Code == "paylater" && statusUser != 2 {
			log.Println("masuk hahaha")
			continue
		}

		if paymentChannel.Code == "paylater" && isPaylater == 0 {
			log.Println("masuk hahaha")
			continue
		}

		if paymentChannel.Code == "tabungan_bima" && statusUser == 0 {
			continue
		}

		paymentChannelResponse.Id = paymentChannel.Id
		paymentChannelResponse.IdPaymentMethod = paymentChannel.IdPaymentMethod
		paymentChannelResponse.Name = paymentChannel.Name
		paymentChannelResponse.Code = paymentChannel.Code
		paymentChannelResponse.Logo = paymentChannel.Logo
		paymentChannelResponse.MethodCode = paymentChannel.PaymentMethod.MethodCode

		// hasilBagi := strconv.FormatFloat(jmlOrder/1000000, 'f', 0, 64)

		if paymentChannel.Code == "paylater" {
			paymentChannelResponse.AdminFee = biayaTanggungRenteng
		} else {
			paymentChannelResponse.AdminFee = paymentChannel.AdminFee
		}

		paymentChannelResponse.AdminFeePercentage = paymentChannel.AdminFeePercentage
		paymentChannelResponses = append(paymentChannelResponses, paymentChannelResponse)
	}
	return paymentChannelResponses
}

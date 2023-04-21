package response

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/ppob"
)

type GetPostpaidTelcoProductResponse struct {
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Status int     `json:"status"`
	Fee    float64 `json:"fee"`
	Komisi float64 `json:"komisi"`
	Type   string  `json:"type"`
}

func ToGetPostpaidTelcoProductResponse(telcoProducts *ppob.InquiryProductPostpaidPPOB) (telcoProductResponses []GetPostpaidTelcoProductResponse) {
	for _, telcoProduct := range telcoProducts.Data.Pasca {
		telcoProductResponse := GetPostpaidTelcoProductResponse{}
		telcoProductResponse.Code = telcoProduct.Code
		telcoProductResponse.Name = telcoProduct.Name
		telcoProductResponse.Status = telcoProduct.Status
		telcoProductResponse.Fee = telcoProduct.Fee
		telcoProductResponse.Komisi = telcoProduct.Komisi
		telcoProductResponse.Type = telcoProduct.Type
		telcoProductResponses = append(telcoProductResponses, telcoProductResponse)
	}
	return telcoProductResponses
}

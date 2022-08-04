package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/ppob"

type GetPostpaidPdamProductResponse struct {
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Status int     `json:"status"`
	Fee    float64 `json:"fee"`
	Komisi float64 `json:"komisi"`
	Type   string  `json:"type"`
}

func ToGetPostpaidPadmProductResponse(pdamProducts []ppob.PostpaidPriceList) (pdamProductResponses []GetPostpaidPdamProductResponse) {
	for _, pdamProduct := range pdamProducts {
		pdamProductResponse := GetPostpaidPdamProductResponse{}
		if pdamProduct.Type == "pdam" {
			pdamProductResponse.Code = pdamProduct.Code
			pdamProductResponse.Name = pdamProduct.Name
			pdamProductResponse.Status = pdamProduct.Status
			pdamProductResponse.Fee = pdamProduct.Fee
			pdamProductResponse.Komisi = pdamProduct.Komisi
			pdamProductResponse.Type = pdamProduct.Type
			pdamProductResponses = append(pdamProductResponses, pdamProductResponse)
		}
	}
	return pdamProductResponses
}

package response

import (
	"strconv"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/ppob"
)

type GetPrepaidPriceListResponse struct {
	ProductCode        string  `json:"product_code"`
	ProductDescription string  `json:"product_description"`
	ProductNominal     string  `json:"product_nominal"`
	ProductDetails     string  `json:"product_details"`
	ProductPrice       float64 `json:"product_price"`
	ProductType        string  `json:"product_type"`
	ActivePeriod       string  `json:"active_period"`
	Status             string  `json:"status"`
	IconUrl            string  `json:"icon_url"`
}

func ToGetPrepaidPriceListResponse(priceLists []ppob.PrepaidPriceList) (priceListResponses []GetPrepaidPriceListResponse) {
	for _, priceList := range priceLists {
		pulsaPriceListResponse := GetPrepaidPriceListResponse{}
		_, err := strconv.Atoi(priceList.ProductNominal)
		if err == nil {
			pulsaPriceListResponse.ProductCode = priceList.ProductCode
			pulsaPriceListResponse.ProductDescription = priceList.ProductDescription
			pulsaPriceListResponse.ProductNominal = priceList.ProductNominal
			pulsaPriceListResponse.ProductDetails = priceList.ProductDetails
			pulsaPriceListResponse.ProductPrice = priceList.ProductPrice + 1500
			pulsaPriceListResponse.ProductType = priceList.ProductType
			pulsaPriceListResponse.ActivePeriod = priceList.ActivePeriod
			pulsaPriceListResponse.Status = priceList.Status
			pulsaPriceListResponse.IconUrl = priceList.IconUrl
			priceListResponses = append(priceListResponses, pulsaPriceListResponse)
		}
	}
	return priceListResponses
}

func ToGetPrepaidDataPriceListResponse(priceLists []ppob.PrepaidPriceList) (priceListResponses []GetPrepaidPriceListResponse) {
	for _, priceList := range priceLists {
		pulsaPriceListResponse := GetPrepaidPriceListResponse{}
		pulsaPriceListResponse.ProductCode = priceList.ProductCode
		pulsaPriceListResponse.ProductDescription = priceList.ProductDescription
		pulsaPriceListResponse.ProductNominal = priceList.ProductNominal
		pulsaPriceListResponse.ProductDetails = priceList.ProductDetails
		pulsaPriceListResponse.ProductPrice = priceList.ProductPrice + 1500
		pulsaPriceListResponse.ProductType = priceList.ProductType
		pulsaPriceListResponse.ActivePeriod = priceList.ActivePeriod
		pulsaPriceListResponse.Status = priceList.Status
		pulsaPriceListResponse.IconUrl = priceList.IconUrl
		priceListResponses = append(priceListResponses, pulsaPriceListResponse)
	}
	return priceListResponses
}

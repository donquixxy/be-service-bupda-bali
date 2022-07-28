package ppob

type PrepaidPriceListResponse struct {
	Data GetPrepaidPriceListData `json:"data"`
}

type GetPrepaidPriceListData struct {
	Data    []PrepaidPriceList `json:"pricelist"`
	Rc      string             `json:"rc"`
	Message string             `json:"message"`
}

type PrepaidPriceList struct {
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

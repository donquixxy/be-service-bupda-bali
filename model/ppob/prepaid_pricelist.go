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

type InquiryPrepaidPln struct {
	Data InquiryPrepaidPlnData `json:"data"`
}

type InquiryPrepaidPlnData struct {
	Status       string `json:"status"`
	CustomerId   string `json:"customer_id"`
	MeterNo      string `json:"meter_no"`
	SubscriberId string `json:"subscriber_id"`
	Name         string `json:"name"`
	SegmentPower string `json:"segment_power"`
	Message      string `json:"message"`
	Rc           string `json:"rc"`
}

type InquiryPrepaidPlnErrorHandle struct {
	Data InquiryPrepaidPlnErrorData `json:"data"`
}

type InquiryPrepaidPlnErrorData struct {
	Rc      string      `json:"rc"`
	Message string      `json:"message"`
	Status  interface{} `json:"status"`
}

type TopupPrepaidResponse struct {
	Data TopupPrepaidResponseData `json:"data"`
}

type TopupPrepaidResponseData struct {
	RefId       string  `json:"ref_id"`
	Status      float64 `json:"status"`
	ProductCode string  `json:"product_code"`
	CustomerId  string  `json:"customer_id"`
	Price       float64 `json:"price"`
	Message     string  `json:"message"`
	Balance     float64 `json:"balance"`
	TrId        int     `json:"tr_id"`
	Rc          string  `json:"rc"`
}

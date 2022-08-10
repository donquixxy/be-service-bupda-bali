package ppob

// Postpaid list
type PostpaidPriceListResponse struct {
	Data GetPostpaidPriceListData `json:"data"`
	Meta []interface{}            `json:"meta"`
}

type GetPostpaidPriceListData struct {
	Pasca []PostpaidPriceList `json:"pasca"`
}

type PostpaidPriceList struct {
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Status int     `json:"status"`
	Fee    float64 `json:"fee"`
	Komisi float64 `json:"komisi"`
	Type   string  `json:"type"`
}

// Inquiry Postpaid PLN
type InquiryPostpaidPln struct {
	Data InquiryPostpaidPlnData `json:"data"`
}

type InquiryPostpaidPlnData struct {
	TrxId        int                    `json:"tr_id"`
	Code         string                 `json:"code"`
	Hp           string                 `json:"hp"`
	TrxName      string                 `json:"tr_name"`
	Period       string                 `json:"period"`
	Nominal      float64                `json:"nominal"`
	Admin        float64                `json:"admin"`
	RefId        string                 `json:"ref_id"`
	ResponseCode string                 `json:"response_code"`
	Message      string                 `json:"message"`
	Price        float64                `json:"price"`
	SellingPrice float64                `json:"selling_price"`
	Desc         InquiryPostpaidPlnDesc `json:"desc"`
}

type InquiryPostpaidPlnDesc struct {
	Tarif         string                    `json:"tarif"`
	Daya          int                       `json:"daya"`
	LembarTagihan string                    `json:"lembar_tagihan"`
	Tagihan       InquiryPostpaidPlnTagihan `json:"tagihan"`
}

type InquiryPostpaidPlnTagihan struct {
	Detail []InquiryPostpaidPlnDetail `json:"detail"`
}

type InquiryPostpaidPlnDetail struct {
	Periode     string  `json:"periode"`
	NilaiTgihan string  `json:"nilai_tagihan"`
	Admin       string  `json:"admin"`
	Denda       string  `json:"denda"`
	Total       float64 `json:"total"`
}

//  Inquiry Postpaid PDAM
type InquiryPostpaidPdam struct {
	Data InquiryPostpaidPdamData `json:"data"`
}

type InquiryPostpaidPdamData struct {
	TrxId        int                     `json:"tr_id"`
	Code         string                  `json:"code"`
	Hp           string                  `json:"hp"` //Customer Number
	TrxName      string                  `json:"tr_name"`
	Period       string                  `json:"period"`
	Nominal      float64                 `json:"nominal"`
	Admin        float64                 `json:"admin"`
	RefId        string                  `json:"ref_id"`
	ResponseCode string                  `json:"response_code"`
	Message      string                  `json:"message"`
	Price        float64                 `json:"price"`
	SellingPrice float64                 `json:"selling_price"`
	Desc         InquiryPostpaidPdamDesc `json:"desc"`
}

type InquiryPostpaidPdamDesc struct {
	BillQuantity int                     `json:"bill_quantity"`
	Address      string                  `json:"address"`
	BillerAdmin  string                  `json:"biller_admin"`
	PdamName     string                  `json:"pdam_name"`
	StampDuty    string                  `json:"stamp_duty"`
	DueDate      string                  `json:"due_date"`
	KodeTarif    string                  `json:"kode_tarif"`
	Bill         InquiryPostpaidPdamBill `json:"bill"`
}

type InquiryPostpaidPdamBill struct {
	Detail []InquiryPostpaidPdamBillDetail `json:"detail"`
}

type InquiryPostpaidPdamBillDetail struct {
	Period     string  `json:"period"`
	FirstMeter int     `json:"first_meter"`
	LastMeter  int     `json:"last_meter"`
	Penalty    float64 `json:"penalty"`
	BillAmount float64 `json:"bill_amount"`
	MiscAmount float64 `json:"misc_amount"`
	Stand      string  `json:"stand"`
}

// Postpaid Check Transaction
type PostpaidCheckTransactionPdam struct {
	Data PostpaidCheckTransactionPdamData `json:"data"`
	Meta []interface{}                    `json:"meta"`
}

type PostpaidCheckTransactionPdamData struct {
	TrxId        int                     `json:"tr_id"`
	Code         string                  `json:"code"`
	Datetime     string                  `json:"datetime"`
	Hp           string                  `json:"hp"`
	TrName       string                  `json:"tr_name"`
	Period       string                  `json:"period"`
	Nominal      float64                 `json:"nominal"`
	Admin        float64                 `json:"admin"`
	Status       int                     `json:"status"`
	ResponseCode string                  `json:"response_code"`
	Message      string                  `json:"message"`
	Price        float64                 `json:"price"`
	SellingPrice float64                 `json:"selling_price"`
	Desc         InquiryPostpaidPdamDesc `json:"desc"`
}

type PostpaidCheckTransactionPln struct {
	Data PostpaidCheckTransactionPlnData `json:"data"`
	Meta []interface{}                   `json:"meta"`
}

type PostpaidCheckTransactionPlnData struct {
	TrxId        int                    `json:"tr_id"`
	Code         string                 `json:"code"`
	Datetime     string                 `json:"datetime"`
	Hp           string                 `json:"hp"`
	TrName       string                 `json:"tr_name"`
	Period       string                 `json:"period"`
	Nominal      float64                `json:"nominal"`
	Admin        float64                `json:"admin"`
	Status       int                    `json:"status"`
	ResponseCode string                 `json:"response_code"`
	Message      string                 `json:"message"`
	Price        float64                `json:"price"`
	SellingPrice float64                `json:"selling_price"`
	Desc         InquiryPostpaidPlnDesc `json:"desc"`
}

type TopupPostaidPdamResponse struct {
	Data int
}

type TopupPostaidPdamDataResponse struct {
	TrxId        int                     `json:"tr_id"`
	Code         string                  `json:"code"`
	Datetime     string                  `json:"datetime"`
	Hp           string                  `json:"hp"` //Customer Number
	TrxName      string                  `json:"tr_name"`
	Period       string                  `json:"period"`
	Nominal      float64                 `json:"nominal"`
	Admin        float64                 `json:"admin"`
	RefId        string                  `json:"ref_id"`
	ResponseCode string                  `json:"response_code"`
	Message      string                  `json:"message"`
	Price        float64                 `json:"price"`
	SellingPrice float64                 `json:"selling_price"`
	Balance      float64                 `json:"balance"`
	NoRef        string                  `json:"noref"`
	Desc         InquiryPostpaidPdamDesc `json:"desc"`
}

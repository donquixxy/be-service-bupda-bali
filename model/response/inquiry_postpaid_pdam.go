package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/ppob"

type InquiryPostpaidPdamResponse struct {
	TrxId        int                                     `json:"trx_id"`
	Code         string                                  `json:"code"`
	Hp           string                                  `json:"hp"`
	TrxName      string                                  `json:"trx_name"`
	Period       string                                  `json:"period"`
	Nominal      float64                                 `json:"nominal"`
	AdminFee     float64                                 `json:"admin_fee"`
	Price        float64                                 `json:"price"`
	SellingPrice float64                                 `json:"selling_price"`
	PdamName     string                                  `json:"pdam_name"`
	RefId        string                                  `json:"ref_id"`
	BillDetail   []InquiryPostpaidPdamBillDetailResponse `json:"bill_detail"`
}

type InquiryPostpaidPdamBillDetailResponse struct {
	Period     string  `json:"period"`
	FirstMeter int     `json:"first_meter"`
	LastMeter  int     `json:"last_meter"`
	Penalty    float64 `json:"penalty"`
	BillAmount float64 `json:"bill_amount"`
	MiscAmount float64 `json:"misc_amount"`
}

func ToInquiryPostpaidPdamResponse(inquiryPostpaidPdam *ppob.InquiryPostpaidPdam, inquiryPostpaidPdamBillDetails []ppob.InquiryPostpaidPdamBillDetail, refId string) (inquiryPostpaidPdamResponse InquiryPostpaidPdamResponse) {
	inquiryPostpaidPdamResponse.TrxId = inquiryPostpaidPdam.Data.TrxId
	inquiryPostpaidPdamResponse.Code = inquiryPostpaidPdam.Data.Code
	inquiryPostpaidPdamResponse.Hp = inquiryPostpaidPdam.Data.Hp
	inquiryPostpaidPdamResponse.TrxName = inquiryPostpaidPdam.Data.TrxName
	inquiryPostpaidPdamResponse.Period = inquiryPostpaidPdam.Data.Period
	inquiryPostpaidPdamResponse.Nominal = inquiryPostpaidPdam.Data.Nominal
	inquiryPostpaidPdamResponse.AdminFee = inquiryPostpaidPdam.Data.Admin
	inquiryPostpaidPdamResponse.Price = inquiryPostpaidPdam.Data.Price
	inquiryPostpaidPdamResponse.SellingPrice = inquiryPostpaidPdam.Data.SellingPrice
	inquiryPostpaidPdamResponse.PdamName = inquiryPostpaidPdam.Data.Desc.PdamName
	inquiryPostpaidPdamResponse.RefId = refId

	var inquiryPostpaidPdamBillDetailResponses []InquiryPostpaidPdamBillDetailResponse
	for _, inquiryPostpaidPdamBillDetail := range inquiryPostpaidPdamBillDetails {
		var inquiryPostpaidPlnBillDetailResponse InquiryPostpaidPdamBillDetailResponse
		inquiryPostpaidPlnBillDetailResponse.Period = inquiryPostpaidPdamBillDetail.Period
		inquiryPostpaidPlnBillDetailResponse.BillAmount = inquiryPostpaidPdamBillDetail.BillAmount
		inquiryPostpaidPlnBillDetailResponse.FirstMeter = inquiryPostpaidPdamBillDetail.FirstMeter
		inquiryPostpaidPlnBillDetailResponse.LastMeter = inquiryPostpaidPdamBillDetail.LastMeter
		inquiryPostpaidPlnBillDetailResponse.Penalty = inquiryPostpaidPdamBillDetail.Penalty
		inquiryPostpaidPlnBillDetailResponse.BillAmount = inquiryPostpaidPdamBillDetail.BillAmount
		inquiryPostpaidPdamBillDetailResponses = append(inquiryPostpaidPdamBillDetailResponses, inquiryPostpaidPlnBillDetailResponse)
	}

	inquiryPostpaidPdamResponse.BillDetail = inquiryPostpaidPdamBillDetailResponses

	return inquiryPostpaidPdamResponse
}

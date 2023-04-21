package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/ppob"

type InquiryPostpaidTelcoResponse struct {
	TrxId        int                                      `json:"trx_id"`
	Code         string                                   `json:"code"`
	Hp           string                                   `json:"hp"`
	TrxName      string                                   `json:"trx_name"`
	Period       string                                   `json:"period"`
	Nominal      float64                                  `json:"nominal"`
	AdminFee     float64                                  `json:"admin_fee"`
	Price        float64                                  `json:"price"`
	SellingPrice float64                                  `json:"selling_price"`
	RefId        string                                   `json:"ref_id"`
	KodeArea     string                                   `json:"kode_area"`
	BillDetail   []InquiryPostpaidTelcoBillDetailResponse `json:"bill_detail"`
}

type InquiryPostpaidTelcoBillDetailResponse struct {
	Period       string `json:"period"`
	NilaiTagihan string `json:"nilai_tagihan"`
	Admin        string `json:"admin"`
	Total        string `json:"total"`
}

func ToInquiryPostpaidTelcoResponse(inquiryPostpaidTelco *ppob.InquiryPostpaidTelco, inquiryPostpaidTelcoBillDetails []ppob.InquiryPostpaidTelcoTagihanData, refId string) (inquiryPostpaidTelcoResponse InquiryPostpaidTelcoResponse) {
	inquiryPostpaidTelcoResponse.Code = inquiryPostpaidTelco.Data.Code
	inquiryPostpaidTelcoResponse.TrxId = inquiryPostpaidTelco.Data.TrxId
	inquiryPostpaidTelcoResponse.Hp = inquiryPostpaidTelco.Data.Hp
	inquiryPostpaidTelcoResponse.TrxName = inquiryPostpaidTelco.Data.TrxName
	inquiryPostpaidTelcoResponse.Period = inquiryPostpaidTelco.Data.Period
	inquiryPostpaidTelcoResponse.Nominal = inquiryPostpaidTelco.Data.Nominal
	inquiryPostpaidTelcoResponse.AdminFee = inquiryPostpaidTelco.Data.Admin
	inquiryPostpaidTelcoResponse.Price = inquiryPostpaidTelco.Data.Price
	inquiryPostpaidTelcoResponse.SellingPrice = inquiryPostpaidTelco.Data.SellingPrice
	inquiryPostpaidTelcoResponse.KodeArea = inquiryPostpaidTelco.Data.Desc.KodeArea
	inquiryPostpaidTelcoResponse.RefId = refId

	var inquiryPostpaidTelcoBillDetailResponses []InquiryPostpaidTelcoBillDetailResponse
	for _, inquiryPostpaidTelcoBillDetail := range inquiryPostpaidTelcoBillDetails {
		var inquiryPostpaidTelcoBillDetailResponse InquiryPostpaidTelcoBillDetailResponse
		inquiryPostpaidTelcoBillDetailResponse.Period = inquiryPostpaidTelcoBillDetail.Periode
		inquiryPostpaidTelcoBillDetailResponse.NilaiTagihan = inquiryPostpaidTelcoBillDetail.NilaiTagihan
		inquiryPostpaidTelcoBillDetailResponse.Admin = inquiryPostpaidTelcoBillDetail.Admin
		inquiryPostpaidTelcoBillDetailResponse.Total = inquiryPostpaidTelcoBillDetail.Total
		inquiryPostpaidTelcoBillDetailResponses = append(inquiryPostpaidTelcoBillDetailResponses, inquiryPostpaidTelcoBillDetailResponse)
	}

	inquiryPostpaidTelcoResponse.BillDetail = inquiryPostpaidTelcoBillDetailResponses

	return inquiryPostpaidTelcoResponse
}

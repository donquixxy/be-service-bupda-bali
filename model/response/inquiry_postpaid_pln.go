package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/ppob"

type InquiryPostpaidPlnResponse struct {
	TrxId        int                                    `json:"trx_id"`
	Code         string                                 `json:"code"`
	TrxName      string                                 `json:"trx_name"`
	Period       string                                 `json:"period"`
	Nominal      float64                                `json:"nominal"`
	Admin        float64                                `json:"admin"`
	TotalBayar   float64                                `json:"total_bayar"`
	SellingPrice float64                                `json:"selling_price"`
	Tarif        string                                 `json:"tarif"`
	Daya         int                                    `json:"daya"`
	RefId        string                                 `json:"ref_id"`
	BillDetail   []InquiryPostpaidPlnBillDetailResponse `json:"bill_detail"`
}

type InquiryPostpaidPlnBillDetailResponse struct {
	Period       string  `json:"period"`
	NilaiTagihan string  `json:"nilai_tagihan"`
	AdminFee     string  `json:"admin_fee"`
	Denda        string  `json:"denda"`
	TotalTagihan float64 `json:"bill_amount"`
}

func ToInquiryPostpaidPlnResponse(inquiryPostpaidPln *ppob.InquiryPostpaidPln, inquiryPostpaidPlnBillDetails []ppob.InquiryPostpaidPlnDetail, refId string) (inquiryPostpaidPlnResponse InquiryPostpaidPlnResponse) {
	inquiryPostpaidPlnResponse.TrxId = inquiryPostpaidPln.Data.TrxId
	inquiryPostpaidPlnResponse.Code = inquiryPostpaidPln.Data.Code
	inquiryPostpaidPlnResponse.TrxName = inquiryPostpaidPln.Data.TrxName
	inquiryPostpaidPlnResponse.Period = inquiryPostpaidPln.Data.Period
	inquiryPostpaidPlnResponse.Nominal = inquiryPostpaidPln.Data.Nominal
	inquiryPostpaidPlnResponse.Admin = inquiryPostpaidPln.Data.Admin
	inquiryPostpaidPlnResponse.TotalBayar = inquiryPostpaidPln.Data.Price
	inquiryPostpaidPlnResponse.SellingPrice = inquiryPostpaidPln.Data.SellingPrice
	inquiryPostpaidPlnResponse.Tarif = inquiryPostpaidPln.Data.Desc.Tarif
	inquiryPostpaidPlnResponse.Daya = inquiryPostpaidPln.Data.Desc.Daya
	inquiryPostpaidPlnResponse.RefId = refId

	var inquiryPostpaidPlnBillDetailResponses []InquiryPostpaidPlnBillDetailResponse
	for _, inquiryPostpaidPlnBillDetail := range inquiryPostpaidPlnBillDetails {
		var inquiryPostpaidPlnBillDetailResponse InquiryPostpaidPlnBillDetailResponse
		inquiryPostpaidPlnBillDetailResponse.Period = inquiryPostpaidPlnBillDetail.Periode
		inquiryPostpaidPlnBillDetailResponse.NilaiTagihan = inquiryPostpaidPlnBillDetail.NilaiTgihan
		inquiryPostpaidPlnBillDetailResponse.AdminFee = inquiryPostpaidPlnBillDetail.Admin
		inquiryPostpaidPlnBillDetailResponse.Denda = inquiryPostpaidPlnBillDetail.Denda
		inquiryPostpaidPlnBillDetailResponse.TotalTagihan = inquiryPostpaidPlnBillDetail.Total
		inquiryPostpaidPlnBillDetailResponses = append(inquiryPostpaidPlnBillDetailResponses, inquiryPostpaidPlnBillDetailResponse)
	}

	inquiryPostpaidPlnResponse.BillDetail = inquiryPostpaidPlnBillDetailResponses

	return inquiryPostpaidPlnResponse
}

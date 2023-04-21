package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"

type FindTunggakanPaylater struct {
	LoanPassdueID          string  `json:"loanPassdueID"`
	LoanPassdueNo          string  `json:"loanPassdueNo"`
	LoanAccountRepaymentID string  `json:"loanAccountRepaymentID"`
	LoanID                 string  `json:"loanID"`
	OverduePrincipal       float64 `json:"overduePrincipal"`
	OverdueInterest        float64 `json:"overdueInterest"`
	OverduePenalty         float64 `json:"overduePenalty"`
	OverdueAmount          float64 `json:"overdueAmount"`
	IsPaid                 bool    `json:"isPaid"`
	IsWaivePenalty         bool    `json:"isWaivePenalty"`
	UserInsert             string  `json:"userInsert"`
	DateInsert             string  `json:"dateInsert"`
	UserUpdate             string  `json:"userUpdate"`
	DateUpdate             string  `json:"dateUpdate"`
	PassdueCode            string  `json:"passdueCode"`
}

func ToFindTunggakanPaylaterResponse(tunggakanPaylater []inveli.TunggakanPaylater) []FindTunggakanPaylater {
	var response []FindTunggakanPaylater
	for _, tunggakan := range tunggakanPaylater {
		response = append(response, FindTunggakanPaylater{
			LoanPassdueID:          tunggakan.LoanPassdueID,
			LoanPassdueNo:          tunggakan.LoanPassdueNo,
			LoanAccountRepaymentID: tunggakan.LoanAccountRepaymentID,
			LoanID:                 tunggakan.LoanID,
			OverduePrincipal:       tunggakan.OverduePrincipal,
			OverdueInterest:        tunggakan.OverdueInterest,
			OverduePenalty:         tunggakan.OverduePenalty,
			OverdueAmount:          tunggakan.OverdueAmount,
			IsPaid:                 tunggakan.IsPaid,
			IsWaivePenalty:         tunggakan.IsWaivePenalty,
			UserInsert:             tunggakan.UserInsert,
			DateInsert:             tunggakan.DateInsert,
			UserUpdate:             tunggakan.UserUpdate,
			DateUpdate:             tunggakan.DateUpdate,
			PassdueCode:            tunggakan.PassdueCode,
		})
	}
	return response
}

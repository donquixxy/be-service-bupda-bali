package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"

type FindTagihanPaylater struct {
	RepaymentDate   string  `json:"repayment_date"`
	RepaymentAmount float64 `json:"repayment_amount"`
	DateInsert      string  `json:"date_insert"`
}

func ToFindTagihanPaylater(riwayatPinjaman []inveli.RiwayatPinjaman2) []FindTagihanPaylater {
	var findTagihanPaylaters []FindTagihanPaylater
	for _, v := range riwayatPinjaman {
		if v.IsPaid {
			continue
		}
		findTagihanPaylater := FindTagihanPaylater{
			RepaymentDate:   v.RepaymentDate,
			RepaymentAmount: v.RepaymentAmount,
			DateInsert:      v.DateInsert,
		}
		findTagihanPaylaters = append(findTagihanPaylaters, findTagihanPaylater)
	}
	return findTagihanPaylaters
}

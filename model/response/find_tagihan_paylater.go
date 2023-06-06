package response

import (
	"time"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"
)

type FindTagihanPaylater struct {
	RepaymentDate   time.Time `json:"repayment_date"`
	RepaymentAmount float64   `json:"repayment_amount"`
	DateInsert      time.Time `json:"date_insert"`
}

type FindTagihanPelunasan struct {
	LoanId          string  `json:"loan_id"`
	RepaymentDate   string  `json:"repayment_date"`
	RepaymentAmount float64 `json:"repayment_amount"`
	DateInsert      string  `json:"date_insert"`
}

type FindDetailPyamentPaylater struct {
	Date        time.Time `json:"date"`
	NoTransaksi string    `json:"no_transaksi"`
	Tagihan     float64   `json:"tagihan"`
	BiayaAdmin  float64   `json:"biaya_admin"`
	Bunga       float64   `json:"bunga"`
	Total       float64   `json:"total"`
}

func ToFindDetailPyamentPaylater(paymentHistory *entity.PaymentHistory) (findPaymentHistory FindDetailPyamentPaylater) {
	findPaymentHistory.Date = paymentHistory.TglPembayaran.Time
	findPaymentHistory.NoTransaksi = paymentHistory.NoTransaksi
	findPaymentHistory.Tagihan = paymentHistory.JmlTagihan
	findPaymentHistory.BiayaAdmin = paymentHistory.BiayaAdmin
	findPaymentHistory.Bunga = paymentHistory.BungaPinjaman
	findPaymentHistory.Total = paymentHistory.Total
	return findPaymentHistory
}

func ToFindTagihanPaylater(order []entity.Order) []FindTagihanPaylater {
	var findTagihanPaylaters []FindTagihanPaylater
	for _, v := range order {
		var findTagihanPaylater FindTagihanPaylater
		// if v.IsPaid {
		// 	continue
		// }

		findTagihanPaylater.RepaymentDate = v.PaymentDueDate.Time
		findTagihanPaylater.RepaymentAmount = v.TotalBill
		findTagihanPaylater.DateInsert = v.OrderedDate

		findTagihanPaylaters = append(findTagihanPaylaters, findTagihanPaylater)

	}

	return findTagihanPaylaters
}

func ToFindTagihanPelunasan(tagihan []inveli.TagihanPaylater) []FindTagihanPelunasan {
	var findTagihanPaylaters []FindTagihanPelunasan
	for _, v := range tagihan {
		var findTagihanPaylater FindTagihanPelunasan

		findTagihanPaylater.LoanId = v.LoanId
		findTagihanPaylater.RepaymentDate = v.EndDate
		findTagihanPaylater.RepaymentAmount = v.LoanAmount
		findTagihanPaylater.DateInsert = v.StartDate

		findTagihanPaylaters = append(findTagihanPaylaters, findTagihanPaylater)

	}

	return findTagihanPaylaters
}

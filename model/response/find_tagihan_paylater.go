package response

import (
	"time"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type FindTagihanPaylater struct {
	RepaymentDate   time.Time `json:"repayment_date"`
	RepaymentAmount float64   `json:"repayment_amount"`
	DateInsert      time.Time `json:"date_insert"`
}

// type TotalTagihan struct {
// 	Total               float64               `json:"total"`
// 	FindTagihanPaylater []FindTagihanPaylater `json:"find_tagihan_paylater"`
// }

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

// func ToFindTagihanPaylater(riwayatPinjaman []inveli.RiwayatPinjaman2) []FindTagihanPaylater {
// 	var findTagihanPaylaters []FindTagihanPaylater
// 	for _, v := range riwayatPinjaman {
// 		var findTagihanPaylater FindTagihanPaylater
// 		// if v.IsPaid {
// 		// 	continue
// 		// }

// 		findTagihanPaylater.RepaymentDate = v.RepaymentDate
// 		findTagihanPaylater.RepaymentAmount = v.RepaymentPrincipal
// 		findTagihanPaylater.DateInsert = v.DateInsert

// 		findTagihanPaylaters = append(findTagihanPaylaters, findTagihanPaylater)

// 	}

// 	return findTagihanPaylaters
// }

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

// func ToFindTunggakanPaylater(tunggakan []inveli.TunggakanPaylater2) []FindTagihanPaylater {
// 	var findTagihanPaylaters []FindTagihanPaylater
// 	for _, v := range tunggakan {
// 		var findTagihanPaylater FindTagihanPaylater
// 		findTagihanPaylater.RepaymentDate = v.DateUpdate
// 		findTagihanPaylater.RepaymentAmount = v.LoanAmount
// 		findTagihanPaylater.DateInsert = v.DateInsert

// 		findTagihanPaylaters = append(findTagihanPaylaters, findTagihanPaylater)
// 	}

// 	return findTagihanPaylaters
// }

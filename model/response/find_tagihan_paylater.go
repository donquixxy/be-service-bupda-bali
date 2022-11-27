package response

import (
	"time"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"
)

type FindTagihanPaylater struct {
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

// func ToFindDetailPyamentPaylater(orders []entity.Order) (findDetailPyamentPaylater FindDetailPyamentPaylater) {

// 	for _, order := range orders {

// 	}
// }

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

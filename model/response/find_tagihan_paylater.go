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

type TotalTagihan struct {
	Total               float64               `json:"total"`
	FindTagihanPaylater []FindTagihanPaylater `json:"find_tagihan_paylater"`
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

func ToFindTagihanPaylater(riwayatPinjaman []inveli.RiwayatPinjaman2) TotalTagihan {
	var findTagihanPaylaters []FindTagihanPaylater
	var total float64
	for _, v := range riwayatPinjaman {
		if v.IsPaid {
			continue
		}
		findTagihanPaylater := FindTagihanPaylater{
			RepaymentDate:   v.RepaymentDate,
			RepaymentAmount: v.RepaymentPrincipal,
			DateInsert:      v.DateInsert,
		}
		total = total + v.RepaymentPrincipal
		findTagihanPaylaters = append(findTagihanPaylaters, findTagihanPaylater)
	}

	totalTagihan := TotalTagihan{
		Total:               total,
		FindTagihanPaylater: findTagihanPaylaters,
	}

	return totalTagihan
}

func ToFindTunggakanPaylater(tunggakan []inveli.TunggakanPaylater) TotalTagihan {
	var findTagihanPaylaters []FindTagihanPaylater
	var total float64
	for _, v := range tunggakan {
		findTagihanPaylater := FindTagihanPaylater{
			RepaymentDate:   v.DateUpdate,
			RepaymentAmount: v.OverdueAmount,
			DateInsert:      v.DateInsert,
		}
		total = total + v.OverdueAmount
		findTagihanPaylaters = append(findTagihanPaylaters, findTagihanPaylater)
	}
	totalTagihan := TotalTagihan{
		Total:               total,
		FindTagihanPaylater: findTagihanPaylaters,
	}
	return totalTagihan
}

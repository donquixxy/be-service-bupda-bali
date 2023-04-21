package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type PaymentHistory struct {
	Id               string    `gorm:"primaryKey;column:id;"`
	NoTransaksi      string    `gorm:"column:no_transaksi;"`
	IdUser           string    `gorm:"column:id_user;"`
	IdDesa           string    `gorm:"column:id_desa;"`
	JmlTagihan       float64   `gorm:"column:jml_tagihan;"`
	BungaPercentage  float64   `gorm:"column:bunga_percentage;"`
	BungaPinjaman    float64   `gorm:"column:bunga_pinjaman;"`
	BiayaAdmin       float64   `gorm:"column:biaya_admin;"`
	Total            float64   `gorm:"column:total;"`
	TglPembayaran    null.Time `gorm:"column:tgl_pembayaran;"`
	StatusPembayaran int       `gorm:"column:status_pembayaran;"`
	IndexDate        string    `gorm:"column:index_date;"`
	CreatedAt        time.Time `gorm:"column:created_at;"`
}

func (PaymentHistory) TableName() string {
	return "payment_history"
}

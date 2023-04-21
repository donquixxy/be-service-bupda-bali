package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type ListPinjaman struct {
	Id               string    `gorm:"primaryKey;column:id;"`
	IdUser           string    `gorm:"column:id_user;"`
	IdDesa           string    `gorm:"column:id_desa;"`
	IdOrder          string    `gorm:"column:id_order;"`
	Nik              string    `gorm:"column:nik;"`
	JmlTagihan       float64   `gorm:"column:jml_tagihan;"`
	BungaPercentage  float64   `gorm:"column:bunga_percentage;"`
	BungaPinjaman    float64   `gorm:"column:bunga_pinjaman;"`
	BiayaAdmin       float64   `gorm:"column:biaya_admin;"`
	Total            float64   `gorm:"column:total;"`
	TglPeminjaman    time.Time `gorm:"column:tgl_peminjaman;"`
	TglJatuhTempo    null.Time `gorm:"column:tgl_jatuh_tempo;"`
	TglPembayaran    null.Time `gorm:"column:tgl_pembayaran;"`
	StatusPembayaran int       `gorm:"column:status_pembayaran;"`
	CreatedAt        time.Time `gorm:"column:created_at;"`
}

func (ListPinjaman) TableName() string {
	return "list_pinjaman"
}

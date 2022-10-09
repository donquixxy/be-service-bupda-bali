package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type MerchantApproveList struct {
	Id            string    `gorm:"primaryKey;column:id;"`
	IdDesa        string    `gorm:"column:id_desa;"`
	IdUser        string    `gorm:"column:id_user;"`
	NamaLengkap   string    `gorm:"column:nama_lengkap;"`
	ApproveStatus int       `gorm:"column:approve_status;"`
	MerchantName  string    `gorm:"column:merchant_name;"`
	CreatedAt     time.Time `gorm:"column:created_at;"`
	UpdatedAt     null.Time `gorm:"column:updated_at;"`
}

func (MerchantApproveList) TableName() string {
	return "merchant_approve_list"
}

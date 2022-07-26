package entity

import "time"

type Promo struct {
	Id          string    `gorm:"primaryKey;column:id;"`
	IdDesa      string    `gorm:"column:id_desa;"`
	PromoTitle  string    `gorm:"column:promo_title;"`
	Description string    `gorm:"column:description;"`
	Image       string    `gorm:"column:image;"`
	StartDate   time.Time `gorm:"column:start_date;"`
	EndDate     time.Time `gorm:"column:end_date;"`
}

func (Promo) TableName() string {
	return "products_promo"
}

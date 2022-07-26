package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Cart struct {
	Id            string       `gorm:"primaryKey;column:id;"`
	IdUser        string       `gorm:"column:id_user;"`
	IdProductDesa string       `gorm:"column:id_product_desa;"`
	ProductsDesa  ProductsDesa `gorm:"foreignKey:IdProductDesa"`
	Qty           int          `gorm:"column:qty;"`
	CreatedDate   time.Time    `gorm:"column:created_at;"`
	UpdatedDate   null.Time    `gorm:"column:updated_at;"`
}

func (Cart) TableName() string {
	return "cart"
}

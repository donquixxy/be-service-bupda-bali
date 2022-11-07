package entity

import (
	"gopkg.in/guregu/null.v4"
)

type UserAccount struct {
	Id          string    `gorm:"primaryKey;column:id;"`
	IdUser      string    `gorm:"column:id_user;"`
	IdAccount   string    `gorm:"column:id_account;"`
	AccountName string    `gorm:"column:account_name;"`
	IdProduct   string    `gorm:"column:id_product;"`
	CreatedDate null.Time `gorm:"column:created_at;"`
	UpdatedDate null.Time `gorm:"column:created_at;"`
}

func (UserAccount) TableName() string {
	return "users_account"
}

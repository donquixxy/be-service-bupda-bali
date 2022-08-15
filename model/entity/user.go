package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type User struct {
	Id              string    `gorm:"primaryKey;column:id;"`
	IdLimitPayLater string    `gorm:"column:id_limit_pay_later;"`
	IdRole          string    `gorm:"column:id_role;"`
	IdDesa          string    `gorm:"column:id_desa;"`
	Desa            Desa      `gorm:"foreignKey:IdDesa;"`
	Phone           string    `gorm:"column:phone;"`
	Password        string    `gorm:"column:password;"`
	RefreshToken    string    `gorm:"column:refresh_token;"`
	LastLoginDate   null.Time `gorm:"column:last_login;"`
	AccountType     int       `gorm:"column:account_type;"`
	StatusSurvey    int       `gorm:"column:status_survey;"`
	MerchantCode    string    `gorm:"column:merchant_code;"`
	TokenDevice     string    `gorm:"column:token_device;"`
	IsActive        int       `gorm:"column:is_active;"`
	IsDelete        int       `gorm:"column:is_delete;"`
	IsDeleteDate    null.Time `gorm:"column:is_delete_date;"`
	CreatedDate     time.Time `gorm:"column:created_at;"`
}

func (User) TableName() string {
	return "users"
}

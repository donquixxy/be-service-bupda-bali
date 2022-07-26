package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type UserProfile struct {
	Id          string    `gorm:"primaryKey;column:id;"`
	IdUser      string    `gorm:"column:id_user;"`
	User        User      `gorm:"foreignKey:IdUser;"`
	NoIdentitas string    `gorm:"column:no_identitas;"`
	NamaLengkap string    `gorm:"column:nama_lengkap;"`
	Email       string    `gorm:"column:email;"`
	CreatedDate time.Time `gorm:"column:created_at;"`
	UpdatedDate null.Time `gorm:"column:created_at;"`
}

func (UserProfile) TableName() string {
	return "users_profile"
}

package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type OtpManager struct {
	Id             string    `gorm:"primaryKey;column:id;"`
	IpAddress      string    `gorm:"column:ip_address;"`
	Phone          string    `gorm:"column:phone;"`
	OtpCode        string    `gorm:"column:otp_code;"`
	OtpExperiedAt  time.Time `gorm:"column:otp_experied_at;"`
	PhoneLimit     int       `gorm:"column:phone_limit;"`
	IpAddressLimit int       `gorm:"column:ip_address_limit;"`
	FreezeDueDate  null.Time `gorm:"column:freeze_due_date;"`
	CreatedDate    time.Time `gorm:"column:created_at;"`
	UpdatedDate    null.Time `gorm:"column:updated_at;"`
}

func (OtpManager) TableName() string {
	return "otp_manager"
}

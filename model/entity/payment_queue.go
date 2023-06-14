package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type PaymentQueue struct {
	Id          string    `gorm:"primaryKey;column:id;"`
	IdUser      string    `gorm:"column:id_user;"`
	IdOrder     string    `gorm:"column:id_order;"`
	LoanId      string    `gorm:"column:loan_id;"`
	Amount      float64   `gorm:"column:amount;"`
	Status      int       `gorm:"column:status;"`
	Order       int       `gorm:"column:order;"`
	CreatedAt   time.Time `gorm:"column:created_at;"`
	ProcessedAt null.Time `gorm:"column:processed_at;"`
}

func (PaymentQueue) TableName() string {
	return "payment_queue"
}

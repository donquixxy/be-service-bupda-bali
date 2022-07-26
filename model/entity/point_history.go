package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type PointHistory struct {
	Id          string    `gorm:"primaryKey;column:id;"`
	IdPoint     string    `gorm:"column:id_point;"`
	IdUser      string    `gorm:"column:id_user;"`
	Debit       float64   `gorm:"column:debit;"`
	Kredit      float64   `gorm:"column:kredit;"`
	Description string    `gorm:"column:description;"`
	TransDate   time.Time `gorm:"column:trans_date;"`
	CreatedDate time.Time `gorm:"column:created_at;"`
	UpdateDate  null.Time `gorm:"column:updated_at;"`
}

func (PointHistory) TableName() string {
	return "Points_history"
}

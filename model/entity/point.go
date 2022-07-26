package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Point struct {
	Id          string    `gorm:"primaryKey;column:id;"`
	IdUser      string    `gorm:"column:id_user;"`
	JmlPoint    float64   `gorm:"column:jml_point;"`
	StatusPoint int       `gorm:"column:status_point;"`
	CreatedDate time.Time `gorm:"column:created_at;"`
	UpdateDate  null.Time `gorm:"column:updated_at;"`
}

func (Point) TableName() string {
	return "points"
}

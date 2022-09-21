package entity

import (
	"gopkg.in/guregu/null.v4"
)

type InfoDesa struct {
	Id          string    `gorm:"primaryKey;column:id;"`
	IdDesa      string    `gorm:"column:id_desa;"`
	Attachments string    `gorm:"column:attachments;"`
	Title       string    `gorm:"column:title;"`
	Content     string    `gorm:"column:content;"`
	IsActive    string    `gorm:"column:is_active;"`
	Url         string    `gorm:"column:url;"`
	CreatedDate null.Time `gorm:"column:created_at;"`
}

func (InfoDesa) TableName() string {
	return "info_desa"
}
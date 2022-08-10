package entity

type InfoDesa struct {
	Id          string `gorm:"primaryKey;column:id;"`
	IdDesa      string `gorm:"column:id_desa;"`
	Attachments string `gorm:"column:attachments;"`
	Title       string `gorm:"column:title;"`
	Content     string `gorm:"column:content;"`
	IsActive    string `gorm:"column:is_active;"`
}

func (InfoDesa) TableName() string {
	return "info_desa"
}

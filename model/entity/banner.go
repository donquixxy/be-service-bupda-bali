package entity

type Banner struct {
	Id              string `gorm:"primaryKey;column:id;"`
	IdDesa          string `gorm:"column:id_desa;"`
	BannerTitle     string `gorm:"column:banner_title;"`
	BannerImg       string `gorm:"column:banner_img;"`
	BannerUrl       string `gorm:"column:banner_url;"`
	BannerReference string `gorm:"column:banner_reference;"`
	IsActive        int    `gorm:"column:is_active;"`
}

func (Banner) TableName() string {
	return "banner"
}

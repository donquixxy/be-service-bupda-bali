package entity

type ProductsPromo struct {
	Id          string `gorm:"primaryKey;column:id;"`
	IdDesa      string `gorm:"column:id_desa;"`
	PromoTitle  string `gorm:"column:promo_title;"`
	Description string `gorm:"column:description;"`
	Images      string `gorm:"column:images;"`
}

func (ProductsPromo) TableName() string {
	return "products_promo"
}

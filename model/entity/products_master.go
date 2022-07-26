package entity

type ProductsMaster struct {
	Id            string  `gorm:"primaryKey;column:id;"`
	IdBrand       int     `gorm:"column:id_brand;"`
	IdCategory    int     `gorm:"column:id_category;"`
	IdSubCategory int     `gorm:"column:id_sub_category;"`
	IdUnit        int     `gorm:"column:id_unit;"`
	NoSku         string  `gorm:"column:no_sku;"`
	ProductName   string  `gorm:"column:product_name;"`
	Price         float64 `gorm:"column:price;"`
	PriceGrosir   float64 `gorm:"column:price_grosir;"`
	Description   string  `gorm:"column:description;"`
	PictureUrl    string  `gorm:"column:picture_url;"`
	Thumbnail     string  `gorm:"column:thumbnail;"`
}

func (ProductsMaster) TableName() string {
	return "products_master"
}

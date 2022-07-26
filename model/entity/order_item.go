package entity

import "time"

type OrderItem struct {
	Id                 string    `gorm:"primaryKey;column:id;"`
	IdOrder            string    `gorm:"column:id_order;"`
	IdProductDesa      string    `gorm:"column:id_product_desa;"`
	IdUser             string    `gorm:"column:id_user;"`
	NumberOrder        string    `gorm:"column:number_order;"`
	NoSku              string    `gorm:"column:no_sku;"`
	ProductName        string    `gorm:"column:product_name;"`
	PictureUrl         string    `gorm:"column:picture_url;"`
	Thumbnail          string    `gorm:"column:thumbnail;"`
	FlagPromo          int       `gorm:"column:flag_promo;"`
	Description        string    `gorm:"column:description;"`
	Weight             float64   `gorm:"column:weight;"`
	Volume             float64   `gorm:"column:volume;"`
	Qty                int       `gorm:"column:qty;"`
	Price              float64   `gorm:"column:price;"`
	PriceAfterDiscount float64   `gorm:"column:price_after_discount;"`
	TotalPrice         float64   `gorm:"column:total_price;"`
	CreatedAt          time.Time `gorm:"column:created_at;"`
}

func (OrderItem) TableName() string {
	return "orders_items"
}

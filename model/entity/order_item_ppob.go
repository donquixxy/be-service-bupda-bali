package entity

import "time"

type OrderItemPpob struct {
	Id           string    `gorm:"primaryKey;column:id;"`
	IdOrder      string    `gorm:"column:id_order;"`
	IdUser       string    `gorm:"column:id_user;"`
	TrId         int       `gorm:"column:tr_id;"`
	RefId        string    `gorm:"column:ref_id;"`
	ProductCode  string    `gorm:"column:product_code;"`
	ProductType  string    `gorm:"column:product_type;"`
	IconUrl      string    `gorm:"column:icon_url;"`
	Nominal      float64   `gorm:"column:nominal;"`
	Admin        float64   `gorm:"column:admin;"`
	TotalTagihan float64   `gorm:"column:total_tagihan;"`
	SellingPrice float64   `gorm:"column:selling_price;"`
	BillDetail   string    `gorm:"column:bill_detail_json;"`
	CreatedAt    time.Time `gorm:"column:created_at;"`
}

func (OrderItemPpob) TableName() string {
	return "orders_items_ppob"
}

package entity

import "time"

type OrderItemPpob struct {
	Id                 string    `gorm:"primaryKey;column:id;"`
	IdOrder            string    `gorm:"column:id_order;"`
	IdUser             string    `gorm:"column:id_user;"`
	ProductCode        string    `gorm:"column:product_code;"`
	ProductName        string    `gorm:"column:product_name;"`
	ProductType        string    `gorm:"column:product_type;"`
	ProductDescription string    `gorm:"column:product_description;"`
	CustomerId         string    `gorm:"column:customer_id;"`
	NoMeter            string    `gorm:"column:no_meter;"`
	SubscriberId       string    `gorm:"column:subscriber_id;"`
	CustomerName       string    `gorm:"column:customer_name;"`
	SegmentPower       string    `gorm:"column:segment_power;"`
	Periode            string    `gorm:"column:periode;"`
	DueDate            string    `gorm:"column:due_date;"`
	RefId              string    `gorm:"column:ref_id;"`
	PdamName           string    `gorm:"column:pdam_name;"`
	PdamAddress        string    `gorm:"column:pdam_address;"`
	PdamKodeTarif      string    `gorm:"column:pdam_kode_tarif;"`
	IconUrl            string    `gorm:"column:icon_url;"`
	Nominal            float64   `gorm:"column:nominal;"`
	Admin              float64   `gorm:"column:admin;"`
	TotalTagihan       float64   `gorm:"column:total_tagihan;"`
	SellingPrice       float64   `gorm:"column:selling_price;"`
	BillDetail         string    `gorm:"column:bill_detail_json;"`
	CreatedAt          time.Time `gorm:"column:created_at;"`
}

func (OrderItemPpob) TableName() string {
	return "orders_items_ppob"
}

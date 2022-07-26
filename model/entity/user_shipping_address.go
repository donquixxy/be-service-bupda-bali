package entity

type UserShippingAddress struct {
	Id               string  `gorm:"primaryKey;column:id;"`
	IdDesa           string  `gorm:"column:id_desa;"`
	IdUser           string  `gorm:"column:id_user;"`
	AlamatPengiriman string  `gorm:"column:alamat_pengiriman;"`
	Latitude         float64 `gorm:"column:latitude;"`
	Longitude        float64 `gorm:"column:longitude;"`
	Radius           float64 `gorm:"column:radius;"`
	StatusPrimary    int     `gorm:"column:is_primary;"`
	Catatan          string  `gorm:"column:catatan;"`
}

func (UserShippingAddress) TableName() string {
	return "users_shipping_address"
}

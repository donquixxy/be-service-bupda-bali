package entity

type PaymentChannel struct {
	Id                 string        `gorm:"primaryKey;column:id;"`
	IdPaymentMethod    string        `gorm:"column:id_payment_method;"`
	PaymentMethod      PaymentMethod `gorm:"foreignKey:IdPaymentMethod"`
	IdDesa             string        `gorm:"column:id_desa;"`
	Name               string        `gorm:"column:name;"`
	Code               string        `gorm:"column:code;"`
	Alias              string        `gorm:"column:alias;"`
	Logo               string        `gorm:"column:logo;"`
	AdminFee           float64       `gorm:"column:admin_fee;"`
	AdminFeePercentage float64       `gorm:"column:admin_fee_percentage;"`
	NoAccountBank      string        `gorm:"column:no_account_bank;"`
	NamaPemilikBank    string        `gorm:"column:nama_pemilik_bank;"`
	IsActive           int           `gorm:"column:is_active;"`
}

func (PaymentChannel) TableName() string {
	return "payment_channel"
}

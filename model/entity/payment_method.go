package entity

type PaymentMethod struct {
	Id         string `gorm:"primaryKey;column:id;"`
	IdDesa     string `gorm:"column:id_desa;"`
	Method     string `gorm:"column:method;"`
	MethodCode string `gorm:"column:method_code;"`
	IsActive   int    `gorm:"column:is_active;"`
}

func (PaymentMethod) TableName() string {
	return "payment_method"
}

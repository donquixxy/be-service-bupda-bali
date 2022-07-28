package entity

type OperatorPrefix struct {
	Id           string `gorm:"primaryKey;column:id;"`
	NamaOperator string `gorm:"column:nama_operator;"`
	KodeOperator string `gorm:"column:kode_operator;"`
	PrefixNumber string `gorm:"column:prefix_number;"`
}

func (OperatorPrefix) TableName() string {
	return "ppob_prefix_operator"
}

package entity

type Kecamatan struct {
	IdKeca   int    `gorm:"primaryKey;column:idkeca;"`
	IdKabu   int    `gorm:"column:idkabu;"`
	IdProp   int    `gorm:"column:idprop;"`
	KdKeca   string `gorm:"column:kdkeca;"`
	NamaKeca string `gorm:"column:nama_keca;"`
}

func (Kecamatan) TableName() string {
	return "master_kecamatan"
}

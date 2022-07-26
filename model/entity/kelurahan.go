package entity

type Kelurahan struct {
	IdKelu   int    `gorm:"primaryKey;column:idkelu;"`
	IdKeca   int    `gorm:"column:idkeca;"`
	IdKabu   int    `gorm:"column:idkabu;"`
	IdProp   int    `gorm:"column:idprop;"`
	KdKelu   string `gorm:"column:kdkelu;"`
	NamaKelu string `gorm:"column:nama_kelu;"`
}

func (Kelurahan) TableName() string {
	return "master_kelurahan"
}

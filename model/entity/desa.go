package entity

type Desa struct {
	Id           string  `gorm:"primaryKey;column:id;"`
	KodeTrx      string  `gorm:"column:kode_trx;"`
	NamaDesa     string  `gorm:"column:nama_desa;"`
	NamaBendesa  string  `gorm:"column:nama_bendesa;"`
	Ongkir       float64 `gorm:"column:ongkir;"`
	GroupIdBupda string  `gorm:"column:group_id;"`
	NoRekening   string  `gorm:"column:no_rekening;"`
}

func (Desa) TableName() string {
	return "desa"
}

package entity

type Desa struct {
	Id          string `gorm:"primaryKey;column:id;"`
	NamaDesa    string `gorm:"column:nama_desa;"`
	NamaBendesa string `gorm:"column:nama_bendesa;"`
}

func (Desa) TableName() string {
	return "desa"
}

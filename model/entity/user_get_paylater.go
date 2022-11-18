package entity

type UserGetPaylater struct {
	Id     string `gorm:"primaryKey;column:id;"`
	NIK    string `gorm:"column:nik;"`
	Nama   string `gorm:"column:nama;"`
	IdDesa string `gorm:"column:id_desa;"`
}

func (UserGetPaylater) TableName() string {
	return "users_get_paylater"
}

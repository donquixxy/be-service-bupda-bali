package entity

type AppVersion struct {
	Id      string `gorm:"primaryKey;column:id;"`
	OS      string `gorm:"column:os;"`
	State   string `gorm:"column:state;"`
	Version string `gorm:"column:ver;"`
}

func (AppVersion) TableName() string {
	return "app_version"
}

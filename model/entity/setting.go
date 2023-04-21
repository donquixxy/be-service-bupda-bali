package entity

type Setting struct {
	Id           int     `gorm:"primaryKey;column:id;"`
	IdDesa       string  `gorm:"column:id_desa;"`
	SettingTitle string  `gorm:"column:settings_title;"`
	SettingName  string  `gorm:"column:settings_name;"`
	Value        float64 `gorm:"column:value;"`
}

func (Setting) TableName() string {
	return "settings"
}

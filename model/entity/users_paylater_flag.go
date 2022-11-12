package entity

import "time"

type UsersPaylaterFlag struct {
	Id                  string    `json:"id"`
	IdUser              string    `json:"id_user"`
	PaylaterDate        time.Time `json:"paylater_date"`
	TanggungRentengFlag int       `json:"tanggung_renteng_flag"`
	CreatedAt           time.Time `json:"created_at"`
}

func (UsersPaylaterFlag) TableName() string {
	return "users_paylater_flag"
}

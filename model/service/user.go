package service

import "time"

type User struct {
	Id           string
	IdDesa       string
	AccountType  int
	Phone        string
	Password     string
	CreatedDate  time.Time
	RefreshToken string
}

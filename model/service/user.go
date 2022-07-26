package service

import "time"

type User struct {
	Id           string
	IdDesa       string
	AccountType  int
	Password     string
	CreatedDate  time.Time
	RefreshToken string
}

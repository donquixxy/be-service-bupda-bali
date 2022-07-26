package service

import "github.com/golang-jwt/jwt"

type Auth struct {
	Id           string `json:"id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenClaims struct {
	Id          string `json:"id"`
	IdDesa      string `json:"id_desa"`
	AccountType int    `json:"account_type"`
	jwt.StandardClaims
}

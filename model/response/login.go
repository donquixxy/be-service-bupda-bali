package response

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func ToLoginResponse(token string, refreshToken string) (loginResponse LoginResponse) {
	loginResponse.Token = token
	loginResponse.RefreshToken = refreshToken
	return loginResponse
}

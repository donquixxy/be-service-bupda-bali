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

type LoginInveliResponse struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"userID"`
}

func ToLoginInveliResponse(accessToken string, userID string) (loginInveliResponse LoginInveliResponse) {
	loginInveliResponse.AccessToken = accessToken
	loginInveliResponse.UserID = userID
	return loginInveliResponse
}

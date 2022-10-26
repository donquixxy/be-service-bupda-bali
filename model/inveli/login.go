package inveli

type InveliLoginModel struct {
	Code         int    `json:"code"`
	AccessToken  string `json:"access_token"`
	UserID       string `json:"userID"`
	Email        string `json:"email"`
	CaptchaToken string `json:"captchaToken"`
}

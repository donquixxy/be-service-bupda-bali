package response

type VerifyOtpResponse struct {
	FormToken string `json:"form_token"`
}

func ToVerifyOtpResponse(token string) (verifyOtpResponse VerifyOtpResponse) {
	verifyOtpResponse.FormToken = token
	return verifyOtpResponse
}

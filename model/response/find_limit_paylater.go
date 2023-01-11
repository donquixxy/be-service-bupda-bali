package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"

type FindLimitPayLaterResponse struct {
	MaxLimit       float64 `json:"max_limit"`
	AvailableLimit float64 `json:"available_limit"`
}

func ToFindLimitPayLaterResponse(limitPayLater *inveli.LimitPaylater, loanAmount float64) (limitPayLaterResponse FindLimitPayLaterResponse) {
	limitPayLaterResponse.MaxLimit = limitPayLater.MaxLimit
	limitPayLaterResponse.AvailableLimit = limitPayLater.MaxLimit - loanAmount
	return limitPayLaterResponse
}

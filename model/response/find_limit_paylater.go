package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"

type FindLimitPayLaterResponse struct {
	MaxLimit       float64 `json:"max_limit"`
	AvailableLimit float64 `json:"available_limit"`
}

func ToFindLimitPayLaterResponse(limitPayLater *inveli.LimitPaylater) (limitPayLaterResponse FindLimitPayLaterResponse) {
	limitPayLaterResponse.MaxLimit = limitPayLater.MaxLimit
	limitPayLaterResponse.AvailableLimit = limitPayLater.AvailableLimit
	return limitPayLaterResponse
}

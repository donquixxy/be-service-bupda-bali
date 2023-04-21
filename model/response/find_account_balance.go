package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"

type FindAccountBalanceResponse struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	ProductName string  `json:"product_name"`
	Balance     float64 `json:"balance"`
}

func ToFindAccountBalanceResponse(userAccount *inveli.InveliAcountInfo) FindAccountBalanceResponse {
	var response FindAccountBalanceResponse
	response.ID = userAccount.ID
	response.Code = userAccount.Code
	response.ProductName = userAccount.ProductName
	response.Balance = userAccount.Balance
	return response
}

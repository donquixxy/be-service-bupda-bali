package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"

type GetMutationResponse struct {
	TransactionDate         string  `json:"transactionDate"`
	TransactionDateCurrency string  `json:"transactionDateCurrency"`
	TransactionType         string  `json:"transactionType"`
	DebitAmount             float64 `json:"debitAmount"`
	CreditAmount            float64 `json:"creditAmount"`
	Description             string  `json:"description"`
}

func ToGetMutationResponse(mutation []inveli.Transaction) (mutationResponse []GetMutationResponse) {
	for _, m := range mutation {
		var response GetMutationResponse
		response.TransactionDate = m.TransactionDate
		response.TransactionDateCurrency = m.TransactionDateCurrency
		response.TransactionType = m.TransactionType
		response.DebitAmount = m.DebitAmount
		response.CreditAmount = m.CreditAmount
		response.Description = m.Description
		mutationResponse = append(mutationResponse, response)
	}
	return mutationResponse
}

package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"

type FindTabunganMobileResponse struct {
	Id          string  `json:"id"`
	ProductName string  `json:"product_name"`
	Balance     float64 `json:"title"`
}

func ToFindTabunganMobileResponse(inveliAccount inveli.InveliAcountInfo) (tabunganResponse FindTabunganMobileResponse) {
	tabunganResponse.Id = inveliAccount.AccountName
	tabunganResponse.ProductName = inveliAccount.ProductName
	tabunganResponse.Balance = inveliAccount.Balance
	return
}

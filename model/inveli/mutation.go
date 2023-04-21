package inveli

type Transaction struct {
	ID                      string  `json:"id"`
	TransactionDate         string  `json:"transactionDate"`
	TransactionDateCurrency string  `json:"transactionDateCurrency"`
	TransactionType         string  `json:"transactionType"`
	DebitAmount             float64 `json:"debitAmount"`
	CreditAmount            float64 `json:"creditAmount"`
	Description             string  `json:"description"`
}

type Data struct {
	Transactions []Transaction `json:"transactions"`
}

type InveliJSONObjectMutation struct {
	Data Data `json:"data"`
}

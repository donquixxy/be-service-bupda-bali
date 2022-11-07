package inveli

type InveliAcountInfo struct {
	ID              string  `json:"id"`
	Code            string  `json:"code"`
	AccountName     string  `json:"accountName"`
	AccountName2    string  `json:"accountName2"`
	RecordStatus    int     `json:"recordStatus"`
	ProductName     string  `json:"productName"`
	ProductID       string  `json:"productID"`
	MemberName      string  `json:"memberName"`
	MemberID        string  `json:"memberID"`
	MemberBranchID  string  `json:"memberBranchID"`
	MemberType      int     `json:"memberType"`
	Email           string  `json:"email"`
	Phone           string  `json:"phone"`
	Balance         float64 `json:"balance"`
	BalanceMerchant int     `json:"balanceMerchant"`
	ClosingBalance  int     `json:"closingBalance"`
	BlockingBalance int     `json:"blockingBalance"`
	ProductType     int     `json:"productType"`
	IsPrimary       bool    `json:"isPrimary"`
}

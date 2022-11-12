package inveli

type TunggakanPaylater struct {
	LoanPassdueID          string  `json:"loanPassdueID"`
	LoanPassdueNo          string  `json:"loanPassdueNo"`
	LoanAccountRepaymentID string  `json:"loanAccountRepaymentID"`
	LoanID                 string  `json:"loanID"`
	OverduePrincipal       float64 `json:"overduePrincipal"`
	OverdueInterest        float64 `json:"overdueInterest"`
	OverduePenalty         float64 `json:"overduePenalty"`
	OverdueAmount          float64 `json:"overdueAmount"`
	IsPaid                 bool    `json:"isPaid"`
	IsWaivePenalty         bool    `json:"isWaivePenalty"`
	UserInsert             string  `json:"userInsert"`
	DateInsert             string  `json:"dateInsert"`
	UserUpdate             string  `json:"userUpdate"`
	DateUpdate             string  `json:"dateUpdate"`
	PassdueCode            string  `json:"passdueCode"`
}

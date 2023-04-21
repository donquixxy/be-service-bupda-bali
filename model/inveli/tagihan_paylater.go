package inveli

type TagihanPaylater struct {
	LoanId     string  `json:"loanID"`
	LoanAmount float64 `json:"loanAmount"`
	StartDate  string  `json:"startDate"`
	EndDate    string  `json:"endDate"`
}

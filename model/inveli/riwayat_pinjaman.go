package inveli

type RiwayatPinjaman struct {
	LoanID               string `json:"loanID"`
	LoanAccountRepayment struct {
		ID                     string        `json:"id"`
		LoanAccountID          string        `json:"loanAccountID"`
		RepaymentType          int           `json:"repaymentType"`
		RepaymentDate          string        `json:"repaymentDate"`
		RepaymentInterest      float64       `json:"repaymentInterest"`
		RepaymentPrincipal     float64       `json:"repaymentPrincipal"`
		RepaymentAmount        float64       `json:"repaymentAmount"`
		RepaymentInterestPaid  float64       `json:"repaymentInterestPaid"`
		RepaymentPrincipalPaid float64       `json:"repaymentPrincipalPaid"`
		OutStandingBakiDebet   float64       `json:"outStandingBakiDebet"`
		TellerID               string        `json:"tellerId"`
		IsPaid                 bool          `json:"isPaid"`
		AmountPaid             float64       `json:"amountPaid"`
		PaymentTxnID           string        `json:"paymentTxnID"`
		RecordStatus           int           `json:"recordStatus"`
		UserInsert             string        `json:"userInsert"`
		DateInsert             string        `json:"dateInsert"`
		UserUpdate             string        `json:"userUpdate"`
		DateUpdate             string        `json:"dateUpdate"`
		LoanPassdues           []interface{} `json:"loanPassdues"`
	} `json:"loanAccountRepayments"`
}

type RiwayatPinjaman2 struct {
	ID                     string  `json:"id"`
	LoanAccountID          string  `json:"loanAccountID"`
	RepaymentType          int     `json:"repaymentType"`
	RepaymentDate          string  `json:"repaymentDate"`
	RepaymentInterest      float64 `json:"repaymentInterest"`
	RepaymentPrincipal     float64 `json:"repaymentPrincipal"`
	RepaymentAmount        float64 `json:"repaymentAmount"`
	RepaymentInterestPaid  float64 `json:"repaymentInterestPaid"`
	RepaymentPrincipalPaid float64 `json:"repaymentPrincipalPaid"`
	OutStandingBakiDebet   float64 `json:"outStandingBakiDebet"`
	TellerID               string  `json:"tellerId"`
	IsPaid                 bool    `json:"isPaid"`
	AmountPaid             float64 `json:"amountPaid"`
	PaymentTxnID           string  `json:"paymentTxnID"`
	RecordStatus           int     `json:"recordStatus"`
	UserInsert             string  `json:"userInsert"`
	DateInsert             string  `json:"dateInsert"`
	UserUpdate             string  `json:"userUpdate"`
	DateUpdate             string  `json:"dateUpdate"`
}

package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/machinebox/graphql"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
)

type InveliTestingController interface {
	GetAccountInfo(c echo.Context) error
	GetStatusAkun(c echo.Context) error
	GetBalanceAccount(c echo.Context) error
	GetRiwayatPinjaman(c echo.Context) error
}

type InveliTestingControllerImplementation struct {
	ConfigurationWebserver config.Webserver
	Logger                 *logrus.Logger
}

func NewInveliTestingController(
	logger *logrus.Logger,
) InveliTestingController {
	return &InveliTestingControllerImplementation{
		Logger: logger,
	}
}

type InveliAcountInfo struct {
	ID              string `json:"id"`
	Code            string `json:"code"`
	AccountName     string `json:"accountName"`
	AccountName2    string `json:"accountName2"`
	RecordStatus    int    `json:"recordStatus"`
	ProductName     string `json:"productName"`
	ProductID       string `json:"productID"`
	MemberName      string `json:"memberName"`
	MemberID        string `json:"memberID"`
	MemberBranchID  string `json:"memberBranchID"`
	MemberType      int    `json:"memberType"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Balance         int    `json:"balance"`
	BalanceMerchant int    `json:"balanceMerchant"`
	ClosingBalance  int    `json:"closingBalance"`
	BlockingBalance int    `json:"blockingBalance"`
	ProductType     int    `json:"productType"`
	IsPrimary       bool   `json:"isPrimary"`
}

type UniversalDTO struct {
	Data interface{} `json:"data"`
}

func (controller *InveliTestingControllerImplementation) GetRiwayatPinjaman(c echo.Context) error {
	IDMember := "62A5AC74-9325-4FA7-8BA0-34E14DDFC808"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGFpbURhdGUiOiIyMDIyLTExLTExVDEzOjQ5OjUxLjUzMyswNzowMCIsImV4cCI6MTY5OTY4NTM5MSwiaWQiOiJOakpCTlVGRE56UXRPVE15TlMwMFJrRTNMVGhDUVRBdE16UkZNVFJFUkVaRE9EQTQiLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoiQ2FyZGxlekFQSSIsInVzZXJFbWFpbCI6InRlbnN1MTA0dGVzdDAxQGdtYWlsLmNvbSIsInVzZXJJRCI6IjYyQTVBQzc0LTkzMjUtNEZBNy04QkEwLTM0RTE0RERGQzgwOCJ9.da4NFKF39nXbsTd_7vB5bHFWOhp-WJABTHbMopO-K3k"

	client := graphql.NewClient(config.GetConfig().Inveli.InveliAPI)

	req := graphql.NewRequest(`
		query ($id: String!) {
			loans(memberID: $id){
        loanID
        code
        customerID
        customerName
        productDesc
        loanProductID
        startDate
        endDate
        tenorMonth
        loanAmount
        interestPercentage
        repaymentMethod
        accountID
        userInsert
        dateInsert
        dateAuthor
        userAuthor
        recordStatus
        isLiquidated
        outstandingAmount
        nominalWajib
        filePDFName
        loanAccountRepayments{
            id
            loanAccountID
            repaymentType
            repaymentDate
            repaymentInterest
            repaymentPrincipal
            repaymentAmount
            repaymentInterestPaid
            repaymentPrincipalPaid
            outStandingBakiDebet
            tellerId
            isPaid
            amountPaid
            paymentTxnID
            recordStatus
            userInsert
            dateInsert
            userUpdate
            dateUpdate
            loanPassdues{
                loanPassdueID
                loanPassdueNo
                loanAccountRepaymentID
                loanID
                overduePrincipal
                overdueInterest
                overduePenalty
                overdueAmount
                isPaid
                isWaivePenalty
                userInsert
                dateInsert
                userUpdate
                dateUpdate
                passdueCode
            }
        }
    	}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("id", IDMember)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return err
	}

	// fmt.Println(respData)

	riwayatPinjamans := []inveli.RiwayatPinjaman2{}
	// var data []interface{}
	for _, loan := range respData.(map[string]interface{})["loans"].([]interface{}) {
		riwayatPinjaman := inveli.RiwayatPinjaman2{}
		riwayatPinjaman.ID = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["id"].(string)
		riwayatPinjaman.LoanAccountID = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["loanAccountID"].(string)
		riwayatPinjaman.RepaymentDate = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["repaymentDate"].(string)
		riwayatPinjaman.RepaymentInterest = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["repaymentInterest"].(float64)
		riwayatPinjaman.RepaymentPrincipal = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["repaymentPrincipal"].(float64)
		riwayatPinjaman.RepaymentAmount = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["repaymentAmount"].(float64)
		riwayatPinjaman.RepaymentInterestPaid = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["repaymentInterestPaid"].(float64)
		riwayatPinjaman.RepaymentPrincipalPaid = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["repaymentPrincipalPaid"].(float64)
		riwayatPinjaman.OutStandingBakiDebet = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["outStandingBakiDebet"].(float64)
		riwayatPinjaman.TellerID = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["tellerId"].(string)
		riwayatPinjaman.IsPaid = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["isPaid"].(bool)
		riwayatPinjaman.AmountPaid = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["amountPaid"].(float64)
		riwayatPinjaman.PaymentTxnID = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["paymentTxnID"].(string)
		riwayatPinjaman.UserInsert = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["userInsert"].(string)
		riwayatPinjaman.DateInsert = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["dateInsert"].(string)
		riwayatPinjaman.UserUpdate = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["userUpdate"].(string)
		riwayatPinjaman.DateUpdate = loan.(map[string]interface{})["loanAccountRepayments"].([]interface{})[0].(map[string]interface{})["dateUpdate"].(string)
		riwayatPinjamans = append(riwayatPinjamans, riwayatPinjaman)
	}

	// fmt.Println("riwayatPinjamans", data)

	responses := response.Response{Code: 201, Mssg: "success", Data: riwayatPinjamans, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *InveliTestingControllerImplementation) GetAccountInfo(c echo.Context) error {

	client := graphql.NewClient("http://api-dev.cardlez.com:8089/query")

	keyword := "af81bd11-b74e-4f27-8a9c-a108588bc410"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGFpbURhdGUiOiIyMDIyLTExLTA5VDE5OjU0OjQxLjQwMSswNzowMCIsImV4cCI6MTY5OTUzNDQ4MSwiaWQiOiJSVFV4TlRSQlJqUXRPREEzUmkwMFJrTkZMVUZCTVVVdFJVUXpSalV3TlVaRk56QTMiLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoiQ2FyZGxlekFQSSIsInVzZXJFbWFpbCI6InRlc3QxODYwNTlAZ21haWwuY29tIiwidXNlcklEIjoiRTUxNTRBRjQtODA3Ri00RkNFLUFBMUUtRUQzRjUwNUZFNzA3In0.x9gVyZYEwaa5AZJDnP7U60Rs4CK2a13sC0SrzPohsSA"

	req := graphql.NewRequest(` 
		query ($keyword: String!) {	
			accounts (search: ID, keyword: $keyword) {
				id
      	code
      	accountName
      	accountName2
      	recordStatus
      	productName
      	productID
      	memberName
      	memberID
     		memberBranchID
    		memberType
    		email
    		phone
    		balance
    		balanceMerchant
    		closingBalance
    		blockingBalance
    		productType
    		isPrimary
			}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("keyword", keyword)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return err
	}
	log.Println("req ", respData)

	userInfos := []InveliAcountInfo{}
	for _, value := range respData.(map[string]interface{})["accounts"].([]interface{}) {
		var userInfo InveliAcountInfo
		userInfo.ID = value.(map[string]interface{})["id"].(string)
		userInfo.Code = value.(map[string]interface{})["code"].(string)
		userInfo.AccountName = value.(map[string]interface{})["accountName"].(string)
		userInfo.AccountName2 = value.(map[string]interface{})["accountName2"].(string)
		userInfo.RecordStatus = int(value.(map[string]interface{})["recordStatus"].(float64))
		userInfo.ProductName = value.(map[string]interface{})["productName"].(string)
		userInfo.ProductID = value.(map[string]interface{})["productID"].(string)
		userInfo.MemberName = value.(map[string]interface{})["memberName"].(string)
		userInfo.MemberID = value.(map[string]interface{})["memberID"].(string)
		userInfo.MemberBranchID = value.(map[string]interface{})["memberBranchID"].(string)
		userInfo.MemberType = int(value.(map[string]interface{})["memberType"].(float64))
		userInfo.Email = value.(map[string]interface{})["email"].(string)
		userInfo.Phone = value.(map[string]interface{})["phone"].(string)
		userInfo.Balance = int(value.(map[string]interface{})["balance"].(float64))
		userInfo.BalanceMerchant = int(value.(map[string]interface{})["balanceMerchant"].(float64))
		userInfo.ClosingBalance = int(value.(map[string]interface{})["closingBalance"].(float64))
		userInfo.BlockingBalance = int(value.(map[string]interface{})["blockingBalance"].(float64))
		userInfo.ProductType = int(value.(map[string]interface{})["productType"].(float64))
		userInfo.IsPrimary = value.(map[string]interface{})["isPrimary"].(bool)
		userInfos = append(userInfos, userInfo)
	}

	fmt.Println(userInfos)

	responses := response.Response{Code: 201, Mssg: "success", Data: userInfos, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *InveliTestingControllerImplementation) GetStatusAkun(c echo.Context) error {
	IDMember := "1793AEF9-ABD4-4333-9CDC-02EDB8C68359"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGFpbURhdGUiOiIyMDIyLTExLTA5VDA3OjI0OjU0LjcxNSswNzowMCIsImV4cCI6MTY5OTQ4OTQ5NCwiaWQiOiJNVGM1TTBGRlJqa3RRVUpFTkMwME16TXpMVGxEUkVNdE1ESkZSRUk0UXpZNE16VTUiLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoiQ2FyZGxlekFQSSIsInVzZXJFbWFpbCI6InRlbnN1MTA0Y3JlYXRvcnM5OEBnbWFpbC5jb20iLCJ1c2VySUQiOiIxNzkzQUVGOS1BQkQ0LTQzMzMtOUNEQy0wMkVEQjhDNjgzNTkifQ.NhlYXCzoJgRJCc4oxLAQIoczsO4x-S_nuCEa3vo5_vc"

	client := graphql.NewClient("http://api-dev.cardlez.com:8089/query")

	req := graphql.NewRequest(` 
		query ($id: String!) {	
			member (id: $id) {
        recordStatus
			}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("id", IDMember)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		// return false, err
	}
	// log.Println("res ", respData.Data.Member.RecordStatus)
	// log.Println("res ", respData.(map[string]interface{})["member"].(map[string]interface{})["recordStatus"].(float64))
	resStatus := respData.(map[string]interface{})["member"].(map[string]interface{})["recordStatus"].(float64)
	if resStatus == 2 {
		responses := response.Response{Code: 201, Mssg: "success", Data: true, Error: []string{}}
		return c.JSON(http.StatusOK, responses)
	} else {
		responses := response.Response{Code: 201, Mssg: "success", Data: false, Error: []string{}}
		return c.JSON(http.StatusOK, responses)
	}
}

func (controller *InveliTestingControllerImplementation) GetBalanceAccount(c echo.Context) error {
	client := graphql.NewClient("http://api-dev.cardlez.com:8089/query")

	IDAccount := "b0b8dd43-4536-41c7-9c2c-fc03048c6c7f"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGFpbURhdGUiOiIyMDIyLTExLTA5VDIyOjIxOjQ1LjYzMCswNzowMCIsImV4cCI6MTY5OTU0MzMwNSwiaWQiOiJSVFV4TlRSQlJqUXRPREEzUmkwMFJrTkZMVUZCTVVVdFJVUXpSalV3TlVaRk56QTMiLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoiQ2FyZGxlekFQSSIsInVzZXJFbWFpbCI6InRlc3QxODYwNTlAZ21haWwuY29tIiwidXNlcklEIjoiRTUxNTRBRjQtODA3Ri00RkNFLUFBMUUtRUQzRjUwNUZFNzA3In0.R0Ng7ZCcW0XrxKlthb2FkXVUIDR2fKhQmiZm-Oq40x4"

	req := graphql.NewRequest(` 
		query ($keyword: String!) {	
			accounts (search: ID, keyword: $keyword) {
				id
      	code
      	accountName
      	accountName2
      	productName
    		balance
    		balanceMerchant
    		closingBalance
    		blockingBalance
			}
		}
	`)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Var("keyword", IDAccount)
	ctx := context.Background()
	var respData interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Println(err)
		return err
	}
	log.Println("req ", respData)

	accountBalance := &inveli.InveliAcountInfo{
		ID:              respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["id"].(string),
		Code:            respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["code"].(string),
		AccountName2:    respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["accountName"].(string),
		ProductName:     respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["productName"].(string),
		Balance:         respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["balance"].(float64),
		BalanceMerchant: int(respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["balanceMerchant"].(float64)),
		ClosingBalance:  int(respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["closingBalance"].(float64)),
		BlockingBalance: int(respData.(map[string]interface{})["accounts"].([]interface{})[0].(map[string]interface{})["blockingBalance"].(float64)),
	}

	fmt.Println("account balance = ", accountBalance.Code)

	return nil

}

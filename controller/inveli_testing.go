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
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
)

type InveliTestingController interface {
	GetAccountInfo(c echo.Context) error
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

func (controller *InveliTestingControllerImplementation) GetAccountInfo(c echo.Context) error {

	client := graphql.NewClient("http://api-dev.cardlez.com:8089/query")

	keyword := "86F4FBA2-3EC7-4ABE-B4E3-3B4851AC9472"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGFpbURhdGUiOiIyMDIyLTExLTA0VDE0OjU4OjA3LjgzOSswNzowMCIsImV4cCI6MTY5OTA4NDY4NywiaWQiOiJPRFpHTkVaQ1FUSXRNMFZETnkwMFFVSkZMVUkwUlRNdE0wSTBPRFV4UVVNNU5EY3kiLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoiQ2FyZGxlekFQSSIsInVzZXJFbWFpbCI6InRlbnN1MTA0dGVzdDAxQGdtYWlsLmNvbSIsInVzZXJJRCI6Ijg2RjRGQkEyLTNFQzctNEFCRS1CNEUzLTNCNDg1MUFDOTQ3MiJ9.0qUPrYUYj6yLx2lA9XOehuKMDOKiRK6qqWXEi_FeMJ8"

	req := graphql.NewRequest(` 
		query ($keyword: String!) {	
			accounts (search: MEMBERID, keyword: $keyword) {
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

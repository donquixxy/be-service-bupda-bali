package controller

import (
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/machinebox/graphql"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
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

func (controller *InveliTestingControllerImplementation) GetAccountInfo(c echo.Context) error {
	client := graphql.NewClient("http://api-dev.cardlez.com:8089/query")

	keyword := "70A8A233-DDED-4552-8613-326F6F4DD2D6"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGFpbURhdGUiOiIyMDIyLTExLTAyVDA5OjMyOjMwLjIwMSswNzowMCIsImV4cCI6MTY5ODg5MjM1MCwiaWQiOiJOekJCT0VFeU16TXRSRVJGUkMwME5UVXlMVGcyTVRNdE16STJSalpHTkVSRU1rUTIiLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoiQ2FyZGxlekFQSSIsInVzZXJFbWFpbCI6InRlbnN1MTA0dGVzdDAxQGdtYWlsLmNvbSIsInVzZXJJRCI6IjcwQThBMjMzLURERUQtNDU1Mi04NjEzLTMyNkY2RjRERDJENiJ9.ZaUGuoZzUCtsVrxH53zrbyGs_44j3sBIw7AZHxarJAk"

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

	fmt.Println(respData)

	return nil
}

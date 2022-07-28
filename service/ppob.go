package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/ppob"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"gorm.io/gorm"
)

type PpobServiceInterface interface {
	GetPrepaidPulsaPriceList(requestId string, numberPhone string) (priceListResponse []response.GetPrepaidPriceListResponse)
	GetPrepaidPriceList(requestId string, id, tipe, operator string) *ppob.PrepaidPriceListResponse
}

type PpobServiceImplementation struct {
	DB                                *gorm.DB
	Logger                            *logrus.Logger
	OperatorPrefixRepositoryInterface repository.OperatorPrefixRepositoryInterface
}

func NewPpobService(
	db *gorm.DB,
	logger *logrus.Logger,
	operatorPrefixRepositoryInterface repository.OperatorPrefixRepositoryInterface,
) PpobServiceInterface {
	return &PpobServiceImplementation{
		DB:                                db,
		Logger:                            logger,
		OperatorPrefixRepositoryInterface: operatorPrefixRepositoryInterface,
	}
}

func (service *PpobServiceImplementation) GetPrepaidPulsaPriceList(requestId string, numberPhone string) (priceListResponses []response.GetPrepaidPriceListResponse) {
	split := strings.Split(numberPhone, "")
	phoneJoin := split[0] + split[1] + split[2] + split[3]

	opereratorPrefixResult, err := service.OperatorPrefixRepositoryInterface.FindOperatorPrefixByPhone(service.DB, phoneJoin)
	exceptions.PanicIfError(err, requestId, service.Logger)

	prepaidPulsaPriceList := service.GetPrepaidPriceList(requestId, numberPhone, "pulsa", opereratorPrefixResult.KodeOperator)
	priceListResponses = response.ToPulsaGetPrepaidPriceListResponse(prepaidPulsaPriceList.Data.Data)
	return priceListResponses
}

func (service *PpobServiceImplementation) GetPrepaidPriceList(requestId string, id, tipe, operator string) *ppob.PrepaidPriceListResponse {

	// Create Request
	body, _ := json.Marshal(map[string]interface{}{
		"status":   "all",
		"username": "087762212544",
		"sign":     "2acaad8d0dd84f9bcff26ea0b0e3af81",
	})

	reqBody := ioutil.NopCloser(strings.NewReader(string(body)))

	urlString := "https://prepaid.iak.dev/api/pricelist?type=" + tipe + "&operator=" + operator
	// URL
	url, _ := url.Parse(urlString)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: reqBody,
	}

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	// Read response body
	// data, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf("body: %s\n", data)

	defer resp.Body.Close()

	prepaidPriceList := &ppob.PrepaidPriceListResponse{}

	if err := json.NewDecoder(resp.Body).Decode(prepaidPriceList); err != nil {
		fmt.Println(err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	if prepaidPriceList.Data.Rc != "00" {
		fmt.Printf("body: %s\n", prepaidPriceList.Data)
	}

	return prepaidPriceList

	// return nil
}

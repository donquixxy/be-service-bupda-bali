package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/payment"
)

type PaymentServiceInterface interface {
	VaQrisPay(requestId string, paymentRequest *payment.IpaymuQrisVaRequest) *payment.IpaymuQrisVaResponse
	CreditCardPay(requestId string, paymentRequest *payment.IpaymuCreditCardRequest) *payment.IpaymuCreditCardResponse
}

type PaymentServiceImplementation struct {
	Logger *logrus.Logger
}

func NewPaymentService(
	logger *logrus.Logger,
) PaymentServiceInterface {
	return &PaymentServiceImplementation{
		Logger: logger,
	}
}

func BodyHash(postBody []byte, ipaymuKey string, ipaymuVa string) (signature string, reqBody io.ReadCloser) {
	bodyHash := sha256.Sum256([]byte(postBody))
	bodyHashToString := hex.EncodeToString(bodyHash[:])
	stringToSign := "POST:" + ipaymuVa + ":" + strings.ToLower(string(bodyHashToString)) + ":" + ipaymuKey

	h := hmac.New(sha256.New, []byte(ipaymuKey))
	h.Write([]byte(stringToSign))
	signature = hex.EncodeToString(h.Sum(nil))

	reqBody = ioutil.NopCloser(strings.NewReader(string(postBody)))

	return signature, reqBody
}

func (service *PaymentServiceImplementation) VaQrisPay(requestId string, paymentRequest *payment.IpaymuQrisVaRequest) *payment.IpaymuQrisVaResponse {
	var ipaymu_va = string(config.GetConfig().IpaymuPayment.IpaymuVa)
	var ipaymu_key = string(config.GetConfig().IpaymuPayment.IpaymuKey)

	url, _ := url.Parse(string(config.GetConfig().IpaymuPayment.IpaymuUrl))

	postBody, _ := json.Marshal(map[string]interface{}{
		"name":           paymentRequest.Name,
		"phone":          paymentRequest.Phone,
		"email":          paymentRequest.Email,
		"amount":         paymentRequest.Amount,
		"notifyUrl":      string(config.GetConfig().IpaymuPayment.IpaymuCallbackUrl),
		"expired":        24,
		"expiredType":    "hours",
		"referenceId":    paymentRequest.ReferenceId,
		"paymentMethod":  paymentRequest.PaymentMethod,
		"paymentChannel": paymentRequest.PaymentChannel,
	})

	signature, reqBody := BodyHash(postBody, ipaymu_key, ipaymu_va)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
			"va":           {ipaymu_va},
			"signature":    {signature},
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
	defer resp.Body.Close()

	var dataResponseIpaymu payment.IpaymuQrisVaResponse

	if err := json.NewDecoder(resp.Body).Decode(&dataResponseIpaymu); err != nil {
		fmt.Println(err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	return &dataResponseIpaymu
}

func (service *PaymentServiceImplementation) CreditCardPay(requestId string, paymentRequest *payment.IpaymuCreditCardRequest) *payment.IpaymuCreditCardResponse {
	var ipaymu_va = string(config.GetConfig().IpaymuPayment.IpaymuVa)
	var ipaymu_key = string(config.GetConfig().IpaymuPayment.IpaymuKey)

	url, _ := url.Parse(string(config.GetConfig().IpaymuPayment.IpaymuSnapUrl	))

	postBody, _ := json.Marshal(map[string]interface{}{
		"product":       paymentRequest.Product,
		"qty":           paymentRequest.Qty,
		"price":         paymentRequest.Price,
		"returnUrl":     string(config.GetConfig().IpaymuPayment.IpaymuThankYouPage),
		"cancelUrl":     string(config.GetConfig().IpaymuPayment.IpaymuCancelUrl),
		"notifyUrl":     string(config.GetConfig().IpaymuPayment.IpaymuCallbackUrl),
		"referenceId":   paymentRequest.ReferenceId,
		"buyerName":     paymentRequest.BuyerName,
		"buyerEmail":    paymentRequest.BuyerEmail,
		"buyerPhone":    paymentRequest.BuyerPhone,
		"paymentMethod": paymentRequest.PaymentMethod,
	})

	bodyHash := sha256.Sum256([]byte(postBody))
	bodyHashToString := hex.EncodeToString(bodyHash[:])
	stringToSign := "POST:" + ipaymu_va + ":" + strings.ToLower(string(bodyHashToString)) + ":" + ipaymu_key

	h := hmac.New(sha256.New, []byte(ipaymu_key))
	h.Write([]byte(stringToSign))
	signature := hex.EncodeToString(h.Sum(nil))

	reqBody := io.NopCloser(strings.NewReader(string(postBody)))

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
			"va":           {ipaymu_va},
			"signature":    {signature},
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
	defer resp.Body.Close()

	var dataResponseIpaymu payment.IpaymuCreditCardResponse

	if err := json.NewDecoder(resp.Body).Decode(&dataResponseIpaymu); err != nil {
		fmt.Println(err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	return &dataResponseIpaymu
}

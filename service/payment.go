package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/payment"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	invelirepository "github.com/tensuqiuwulu/be-service-bupda-bali/repository/inveli_repository"
	"gorm.io/gorm"
)

type PaymentServiceInterface interface {
	VaQrisPay(requestId string, paymentRequest *payment.IpaymuQrisVaRequest) *payment.IpaymuQrisVaResponse
	CreditCardPay(requestId string, paymentRequest *payment.IpaymuCreditCardRequest) *payment.IpaymuCreditCardResponse
	CheckPaymentStatus(requestId string, trxId int) *payment.PaymentStatus
	PayPaylater(requestId, idUser string) error
	GetRiwayatPaylaterPerBulan(requestId, idUser string) (response []response.GetRiwayatPaylaterPerbulanResponse)
}

type PaymentServiceImplementation struct {
	DB                           *gorm.DB
	Logger                       *logrus.Logger
	UserRepositoryInterface      repository.UserRepositoryInterface
	InveliAPIRepositoryInterface invelirepository.InveliAPIRepositoryInterface
	OrderRepositoryInterface     repository.OrderRepositoryInterface
}

func NewPaymentService(
	db *gorm.DB,
	logger *logrus.Logger,
	userRepository repository.UserRepositoryInterface,
	inveliAPIRepository invelirepository.InveliAPIRepositoryInterface,
	orderRepository repository.OrderRepositoryInterface,
) PaymentServiceInterface {
	return &PaymentServiceImplementation{
		DB:                           db,
		Logger:                       logger,
		UserRepositoryInterface:      userRepository,
		InveliAPIRepositoryInterface: inveliAPIRepository,
		OrderRepositoryInterface:     orderRepository,
	}
}

func (service *PaymentServiceImplementation) GetRiwayatPaylaterPerBulan(requestId string, idUser string) (response []response.GetRiwayatPaylaterPerbulanResponse) {

	return nil
}

func (service *PaymentServiceImplementation) PayPaylater(requestId, idUser string) error {
	var err error

	user, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)

	if err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	// cek tagihan
	tagihan, err := service.InveliAPIRepositoryInterface.GetTagihanPaylater(user.User.InveliIDMember, user.User.InveliAccessToken)

	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger)
	}

	if tagihan == nil {
		exceptions.PanicIfRecordNotFound(errors.New("TAGIHAN NOT FOUND"), requestId, []string{"TAGIHAN NOT FOUND"}, service.Logger)
	}

	err = service.InveliAPIRepositoryInterface.PayPaylater(user.User.InveliIDMember, user.User.InveliAccessToken)
	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger)
	}

	service.OrderRepositoryInterface.UpdateOrderPaylaterPaidStatus(service.DB, idUser, &entity.Order{
		PaylaterPaidStatus: 1,
	})

	return nil
}

func (service *PaymentServiceImplementation) CheckPaymentStatus(requestId string, trxId int) *payment.PaymentStatus {
	var ipaymu_va = string(config.GetConfig().IpaymuPayment.IpaymuVa)
	var ipaymu_key = string(config.GetConfig().IpaymuPayment.IpaymuKey)

	url, _ := url.Parse(string(config.GetConfig().IpaymuPayment.IpaymuTranscationUrl))

	postBody, _ := json.Marshal(map[string]interface{}{
		"transactionId": strconv.Itoa(trxId),
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

	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	dataResponseIpaymu := &payment.PaymentStatusResponse{}

	if err = json.Unmarshal([]byte(data), &dataResponseIpaymu); err != nil {
		exceptions.PanicIfBadRequest(errors.New("INVALID DATA"), requestId, []string{"INVALID DATA"}, service.Logger)
	}

	return &dataResponseIpaymu.Data
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

	url, _ := url.Parse(string(config.GetConfig().IpaymuPayment.IpaymuSnapUrl))

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

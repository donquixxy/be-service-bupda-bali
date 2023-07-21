package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/payment"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	invelirepository "github.com/tensuqiuwulu/be-service-bupda-bali/repository/inveli_repository"
	"github.com/tensuqiuwulu/be-service-bupda-bali/utilities"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type PaymentServiceInterface interface {
	VaQrisPay(requestId string, paymentRequest *payment.IpaymuQrisVaRequest) *payment.IpaymuQrisVaResponse
	CreditCardPay(requestId string, paymentRequest *payment.IpaymuCreditCardRequest) *payment.IpaymuCreditCardResponse
	CheckPaymentStatus(requestId string, trxId int) *payment.PaymentStatus
	PayPaylaterBill(requestId, idUser string) error
	GetTabunganBimaMutation(requestId, idUser, startDate, endDate string) (responseMutation []response.GetMutationResponse)
	PayWithPaylater(inveliAccessToken, inveliIdMember, desaGroupIdBupda, desaNoRekening, idUser string, orderRequestTotalBill, orderRequestPaymentFee float64) entity.Order
	DebetMultipleTransaksi(requestId, idUser string, loanId []string) error
	GetTagihanPelunasan(requestId string, idUser string) (tagihanPaylaterResponse []response.FindTagihanPelunasan)
}

type PaymentServiceImplementation struct {
	DB                                *gorm.DB
	Logger                            *logrus.Logger
	UserRepositoryInterface           repository.UserRepositoryInterface
	InveliAPIRepositoryInterface      invelirepository.InveliAPIRepositoryInterface
	OrderRepositoryInterface          repository.OrderRepositoryInterface
	PaymentHistoryRepositoryInterface repository.PaymentHistoryRepositoryInterface
	PaymentQueueRepositoryInterface   repository.PaymentQueueRepositoryInterface
}

func NewPaymentService(
	db *gorm.DB,
	logger *logrus.Logger,
	userRepository repository.UserRepositoryInterface,
	inveliAPIRepository invelirepository.InveliAPIRepositoryInterface,
	orderRepository repository.OrderRepositoryInterface,
	paymentHistoryRepository repository.PaymentHistoryRepositoryInterface,
	paymentQueueRepository repository.PaymentQueueRepositoryInterface,
) PaymentServiceInterface {
	return &PaymentServiceImplementation{
		DB:                                db,
		Logger:                            logger,
		UserRepositoryInterface:           userRepository,
		InveliAPIRepositoryInterface:      inveliAPIRepository,
		OrderRepositoryInterface:          orderRepository,
		PaymentHistoryRepositoryInterface: paymentHistoryRepository,
		PaymentQueueRepositoryInterface:   paymentQueueRepository,
	}
}

// Digunakan untuk memproses list tunggakan yang ada.
// Frontnend memilih transaksi2 mana yang ingin dibayar
// Karena pihak inveli tidak bisa melakukan multiple pelunasan.
// Maka harus kita loop request (ID Loan) yang dikirim oleh frontend
// Kemudian HIT API Inveli untuk melakukan pelunasan
func (service *PaymentServiceImplementation) DebetMultipleTransaksi(requestId, idUser string, loanId []string) error {
	jmlLoan := len(loanId)
	order := 0
	orderOldest, err := service.OrderRepositoryInterface.FindOldestUnPaidPaylater(service.DB, idUser, jmlLoan)
	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{"order not found"}, service.Logger)
	}

	paymentQueues := make([]entity.PaymentQueue, len(loanId))
	for i, loan := range loanId {
		order++
		paymentQueue := entity.PaymentQueue{
			Id:          utilities.RandomUUID(),
			IdOrder:     orderOldest[i].Id,
			IdUser:      idUser,
			LoanId:      loan,
			Amount:      orderOldest[i].TotalBill,
			Status:      0,
			Order:       order,
			CreatedAt:   time.Now(),
			ProcessedAt: null.TimeFrom(time.Now()),
		}
		paymentQueues[i] = paymentQueue
	}

	if len(paymentQueues) <= 0 {
		exceptions.PanicIfBadRequest(errors.New("empty payment queues"), requestId, []string{"error creating payment queue: empty payment queues"}, service.Logger)
	}

	err = service.PaymentQueueRepositoryInterface.CreatePaymentQueue(service.DB, paymentQueues)
	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{"error creating payment queue: " + err.Error()}, service.Logger)
	}

	go service.QueuePaymentProcces(idUser)

	return nil
}

func (service *PaymentServiceImplementation) QueuePaymentProcces(idUser string) {
	var err error
	test := 0
	user, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	if err != nil {
		exceptions.PanicIfBadRequest(err, "", []string{"user not found"}, service.Logger)
	}

	paymentQueues, _ := service.PaymentQueueRepositoryInterface.FindPaymentQueueByIdUser(service.DB, idUser)

	log.Println("paymentQueues", len(paymentQueues))
	exceptions.PanicIfBadRequest(err, "", []string{"user not found"}, service.Logger)

	for _, paymentQueue := range paymentQueues {
		test++
		log.Println("count payment queue", test)
		err = service.InveliAPIRepositoryInterface.DebetPerTransaksi(user.User.InveliAccessToken, paymentQueue.LoanId)
		if err != nil {
			log.Println("error debet per transaksi", err.Error())
			err = service.PaymentQueueRepositoryInterface.UpdateFailedPaymentQueueById(service.DB, idUser, &entity.PaymentQueue{
				Status:      2,
				ProcessedAt: null.TimeFrom(time.Now()),
			})

			if err != nil {
				exceptions.PanicIfBadRequest(err, "", []string{"error update payment queue: " + err.Error()}, service.Logger)
			}

			break
		}

		err = service.OrderRepositoryInterface.UpdateOrderPaylaterPaidStatusByIdOrder(service.DB, paymentQueue.IdOrder, &entity.Order{
			PaylaterPaidStatus: 1,
		})

		if err != nil {
			exceptions.PanicIfBadRequest(err, "", []string{"error update order paylater status: " + err.Error()}, service.Logger)
		}

		service.PaymentQueueRepositoryInterface.UpdatePaymentQueueById(service.DB, paymentQueue.Id, &entity.PaymentQueue{
			Status:      1,
			ProcessedAt: null.TimeFrom(time.Now()),
		})

		if err != nil {
			exceptions.PanicIfBadRequest(err, "", []string{"error update payment queue: " + err.Error()}, service.Logger)
		}

	}
}

func (service *PaymentServiceImplementation) PayWithPaylater(inveliAccessToken, inveliIdMember, desaGroupIdBupda, desaRekening, idUser string, orderRequestTotalBill, orderRequestPaymentFee float64) entity.Order {
	var isMerchant float64
	var totalAmount float64
	var err error

	// cek tunggakan 2 bulan terakhir
	unPaidPaylater, err := service.OrderRepositoryInterface.FindUnPaidPaylater(service.DB, idUser)
	if err != nil {
		exceptions.PanicIfError(err, "", service.Logger)
	}

	if len(unPaidPaylater) > 0 && unPaidPaylater[0].OrderedDate.Before(time.Now().AddDate(0, -2, 0)) {
		log.Println("ada tunggakan !!! ")
		exceptions.PanicIfBadRequest(errors.New("masih ada tunggakan 2 bulan terakhir"), "", []string{"Masih ada tunggakan selama 2 bulan yang belum dibayar"}, service.Logger)
	}

	// Set Is Merchant 0
	isMerchant = 0

	totalAmount = orderRequestTotalBill + orderRequestPaymentFee

	// Validasi Saldo Bupda
	saldoBupda, err := service.InveliAPIRepositoryInterface.GetSaldoBupda(inveliAccessToken, desaGroupIdBupda)

	if err != nil {
		exceptions.PanicIfRecordNotFound(errors.New("error saldo bupda "+err.Error()), "", []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger)
	}

	if saldoBupda <= 0 {
		exceptions.PanicIfRecordNotFound(errors.New("saldo bupda kurang"), "", []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger)
	}

	// Get Bunga
	bunga, errr := service.InveliAPIRepositoryInterface.GetLoanProduct(inveliAccessToken)
	if errr != nil {
		exceptions.PanicIfRecordNotFound(errors.New("error get loan product "+err.Error()), "", []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger)
	}

	// Get Loan Product
	loandProductID, err := service.InveliAPIRepositoryInterface.GetLoanProductId(inveliAccessToken)
	if errr != nil {
		exceptions.PanicIfRecordNotFound(errors.New("error get loan product id "+err.Error()), "", []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger)
	}

	if len(loandProductID) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("loan product id not found"), "", []string{"loan product id not found"}, service.Logger)
	}

	// Get Account User
	accountUser, _ := service.UserRepositoryInterface.GetUserAccountPaylaterByID(service.DB, idUser)
	if len(accountUser.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user account paylater not found"), "", []string{"user account paylater not found"}, service.Logger)
	}

	// Validasi Tunggakan Paylater
	// logic untuk validasi tunggakan

	err = service.InveliAPIRepositoryInterface.InveliCreatePaylater(inveliAccessToken, inveliIdMember, accountUser.IdAccount, orderRequestTotalBill, totalAmount, isMerchant, bunga, loandProductID, desaRekening)
	if err != nil {
		exceptions.PanicIfRecordNotFound(errors.New("error care pinjaman "+err.Error()), "", []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger)
	}

	orderEntity := entity.Order{}

	orderEntity.PaymentDueDate = null.NewTime(time.Now().AddDate(0, 0, 30), true)

	orderEntity.OrderStatus = 1
	orderEntity.PaymentStatus = 1
	orderEntity.PaymentName = "Paylater"
	orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
	orderEntity.PaymentCash = totalAmount

	var jmlOrder float64

	// Cek belanjaan user (MENGGUNAKAN PAYLATER) pada periode ke X
	// Selanjutnya kita sebut sebagai jmlOrderPayLate
	jmlOrderPayLate, err := service.OrderRepositoryInterface.FindOrderPayLaterById(service.DB, idUser)
	if err != nil {
		log.Println(err.Error())
	}
	jmlOrder = 0
	for _, v := range jmlOrderPayLate {
		fmt.Println("Order total bill :", v.TotalBill)
		jmlOrder = jmlOrder + v.TotalBill
		fmt.Println("Jml order :", jmlOrder)
	}

	// Get flagging tanggung renteng user pada periode ke X
	userPaylaterFlag, _ := service.UserRepositoryInterface.GetUserPayLaterFlagThisMonth(service.DB, idUser)

	// jmlOrderPayLate ditambah dengan request total bill
	// Flagging tanggung renteng user pada periode ke X dikali 1JUTA

	// Contoh case
	// User A punya limit Paylater 500ribu
	// User A belanja item X seharga 200ribu
	// Maka, user harus membayar 202.500 Ribu. *USER MASIH MEMILKI SISA LIMIT 300RIBU*
	// Kenapa ? karena user pada periode tersebut melakukan transaksi menggunakan paylater
	// Kemudian, user melakukan transaksi sebesar 200ribu *USER MASIH MEMILIKI SISA LIMIT 100RIBU*
	// Maka user hanya membayar 200ribu. KENAPA ?
	// Karena user telah membayar 2500 (Tanggung renteng pada periode yang sama)

	// Case 2
	// User C punya Limit paylater 500ribu
	// User C belanja item X seharga 500ribu
	// Maka, user harus membayar sebesar 502.500 ribu *SISA LIMIT 0*
	// Masih dalam bulan (periode) yang sama. User ternyata melunasi TAGIHAN PAYLATER YANG IA MILIKI (500ribu)
	// Maka, limit paylater user akan kembali ke 500ribu
	// Kenapa ? karena user telah melunasi tagihan pada periode sebelumnya
	// Kemudian, user melakukan transaksi 500ribu
	// Maka user harus membayar sebesar 502.500 lagi
	// Kenapa 2500 ? kenapa tidak 5000 ?
	// Karena setelah user melakukan pembayaran tagihan
	// Bisa dibilang user mendapatkan *reset* untuk tanggung renteng nya

	// Cek jika total Belanjaan user pada periode ke X + dengan order request yang masuk
	// Lebih dari 1 juta

	if (int(jmlOrder) + int(orderRequestTotalBill)) > (userPaylaterFlag.TanggungRentengFlag * 1000000) {
		// Maka update paylaterflag nya
		service.UserRepositoryInterface.UpdateUserPayLaterFlag(service.DB, idUser, &entity.UsersPaylaterFlag{
			TanggungRentengFlag: userPaylaterFlag.TanggungRentengFlag + 1,
		})
	}

	return orderEntity
}

func (service *PaymentServiceImplementation) GetTabunganBimaMutation(requestId, idUser, startDate, endDate string) (responseMutation []response.GetMutationResponse) {
	var err error

	user, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	if err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	if len(user.User.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user not found"), requestId, []string{"user not found"}, service.Logger)
	}

	account, err := service.UserRepositoryInterface.GetUserAccountBimaByID(service.DB, idUser)
	if err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	if len(account.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("account not found"), requestId, []string{"account not found"}, service.Logger)
	}

	mutation, err := service.InveliAPIRepositoryInterface.GetMutation(user.User.InveliAccessToken, account.IdAccount, startDate, endDate)
	if err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	if len(mutation) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("mutation not found"), requestId, []string{"mutation not found"}, service.Logger)
	}

	log.Println("mutation", mutation)

	response := response.ToGetMutationResponse(mutation)
	return response
}

func (service *PaymentServiceImplementation) PayPaylaterBill(requestId, idUser string) error {
	var err error

	user, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)

	if err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}
	// cek tagihan
	tagihan, err := service.InveliAPIRepositoryInterface.GetTagihanPaylater(user.User.InveliIDMember, user.User.InveliAccessToken)

	for _, item := range tagihan {
		fmt.Println("Amount paid :", item.AmountPaid)
		fmt.Println("Is paid :", item.IsPaid)
		fmt.Println("Repayment amount :", item.RepaymentAmount)
	}

	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger)
	}

	if tagihan == nil {
		exceptions.PanicIfRecordNotFound(errors.New("TAGIHAN NOT FOUND"), requestId, []string{"TAGIHAN NOT FOUND"}, service.Logger)
	}

	// get order
	var totalTagihan float64
	var totalBill float64
	var adminFee float64
	var subTotal float64
	order, _ := service.OrderRepositoryInterface.FindOrderPaylaterUnpaidById(service.DB, idUser)
	for _, v := range order {
		totalTagihan += v.PaymentCash
		adminFee += v.PaymentFee
		totalBill += v.TotalBill
		subTotal += v.SubTotal
	}

	err = service.InveliAPIRepositoryInterface.PayPaylater(user.User.InveliIDMember, user.User.InveliAccessToken)
	if err != nil {
		exceptions.PanicIfRecordNotFound(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger)
	}

	now := time.Now()
	generateCode := 100000 + rand.Intn(999999-100000)
	numberOrder := "TAGIHAN" + "/" + now.Format("20060102") + "/" + fmt.Sprint(generateCode)

	orderEntity := &entity.Order{}
	orderEntity.Id = utilities.RandomUUID()
	orderEntity.IdUser = idUser
	orderEntity.NumberOrder = numberOrder
	orderEntity.ProductType = "payment"
	orderEntity.OrderType = 9
	orderEntity.NamaLengkap = user.NamaLengkap
	orderEntity.Email = user.Email
	orderEntity.Phone = user.User.Phone
	orderEntity.ShippingCost = 0
	orderEntity.PaymentCash = totalTagihan
	orderEntity.PaymentFee = adminFee
	orderEntity.SubTotal = subTotal
	orderEntity.TotalBill = totalBill
	orderEntity.PaymentMethod = "tabungan_bima"
	orderEntity.PaymentChannel = "tabungan_bima"
	orderEntity.PaymentName = "Tabungan Bima"
	orderEntity.OrderStatus = 5
	orderEntity.PaymentStatus = 1
	orderEntity.OrderedDate = time.Now()
	orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)

	tx := service.DB.Begin()

	err = service.OrderRepositoryInterface.CreateOrder(tx, orderEntity)
	if err != nil {
		exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error insert payment history " + err.Error()}, service.Logger, tx)
	}

	err = service.OrderRepositoryInterface.UpdateOrderPaylaterPaidStatus(tx, idUser, &entity.Order{
		PaylaterPaidStatus: 1,
	})
	if err != nil {
		exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error update paylater status " + err.Error()}, service.Logger, tx)
	}

	commit := tx.Commit()
	exceptions.PanicIfError(commit.Error, requestId, service.Logger)

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

	data, _ := io.ReadAll(resp.Body)
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

	reqBody = io.NopCloser(strings.NewReader(string(postBody)))

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

func (service *PaymentServiceImplementation) GetTagihanPelunasan(requestId string, idUser string) (tagihanPaylaterResponse []response.FindTagihanPelunasan) {
	user, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{"user not found"}, service.Logger)
	}

	tagihanPaylater, err := service.InveliAPIRepositoryInterface.GetTagihanPaylaterByLatest(user.User.InveliIDMember, user.User.InveliAccessToken)
	if err != nil {
		log.Println("error get tagihan inveli", err.Error())
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	tagihanPaylaterResponse = response.ToFindTagihanPelunasan(tagihanPaylater)
	return tagihanPaylaterResponse

}

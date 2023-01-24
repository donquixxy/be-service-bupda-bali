package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime"

	"log"
	"math"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/payment"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/ppob"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	invelirepository "github.com/tensuqiuwulu/be-service-bupda-bali/repository/inveli_repository"
	"github.com/tensuqiuwulu/be-service-bupda-bali/utilities"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type OrderServiceInterface interface {
	CreateOrderSembako(requestId, idUser, idDesa string, accountType int, orderRequest *request.CreateOrderRequest) (createOrderResponse response.CreateOrderResponse)
	CreateOrderPrepaidPulsa(requestId, idUser, idDesa, productType string, orderRequest *request.CreateOrderPrepaidRequest) (createOrderResponse response.CreateOrderResponse)
	CreateOrderPrepaidPln(requestId, idUser, idDesa, productType string, orderRequest *request.CreateOrderPrepaidRequest) (createOrderResponse response.CreateOrderResponse)
	CreateOrderPostpaidPdam(requestId, idUser, idDesa, productType string, orderRequest *request.CreateOrderPostpaidRequest) (createOrderResponse response.CreateOrderResponse)
	CreateOrderPostpaidTelco(requestId, idUser, idDesa, productType string, orderRequest *request.CreateOrderPostpaidRequest) (createOrderResponse response.CreateOrderResponse)
	CreateOrderPostpaidPln(requestId, idUser, idDesa, productType string, orderRequest *request.CreateOrderPostpaidRequest) (createOrderResponse response.CreateOrderResponse)
	FindOrderByUser(requestId, idUser string, orderStatus int) (orderResponses []response.FindOrderByUserResponse)
	FindOrderSembakoById(requestId, idOrder string) (orderResponse response.FindOrderSembakoByIdResponse)
	FindOrderPrepaidPulsaById(requestId, idOrder string, productType string) (orderResponse response.FindOrderPrepaidPulsaByIdResponse)
	FindOrderPrepaidPlnById(requestId, idOrder string) (orderResponse response.FindOrderPrepaidPlnByIdResponse)
	FindOrderPostpaidPdamById(requestId, idOrder string) (orderResponse response.FindOrderPostpaidPdamByIdResponse)
	FindOrderPostpaidPlnById(requestId, idOrder string) (orderResponse response.FindOrderPostpaidPlnByIdResponse)
	CancelOrderById(requestId string, orderRequest *request.OrderIdRequest)
	CompleteOrderById(requestId string, orderRequest *request.OrderIdRequest)
	UpdatePaymentStatusOrder(requestId string, orderRequest *request.UpdatePaymentStatusOrderRequest)
	GenerateNumberOrder(idDesa string) (numberOrder string)
	PrepaidPulsaTopup(requestId string, customerId, refId, productCode string) *ppob.TopupPrepaidPulsaResponse
	OrderInquiryPrepaidPln(requestId string, customerId string) (inquiryPrepaidPlnResponse response.InquiryPrepaidPlnResponse)
	CallbackPpobTransaction(requestId string, ppobCallbackRequest *request.PpobCallbackRequest)
	FindOrderPayLaterByIdUser(requestId, idUser string) (orderResponse []response.FindOrderByUserResponse)
	FindOrderPaymentById(requestId, idOrder string) (orderResponse response.OrderPayment)
	SendMessageToTelegram(message, chatId, token string)
}

type OrderServiceImplementation struct {
	DB                                     *gorm.DB
	Validate                               *validator.Validate
	Logger                                 *logrus.Logger
	OrderRepositoryInterface               repository.OrderRepositoryInterface
	UserRepositoryInterface                repository.UserRepositoryInterface
	PaymentServiceInterface                PaymentServiceInterface
	CartRepositoryInterface                repository.CartRepositoryInterface
	OrderItemRepositoryInterface           repository.OrderItemRepositoryInterface
	PaymentChannelRepositoryInterface      repository.PaymentChannelRepositoryInterface
	ProductDesaRepositoryInterface         repository.ProductDesaRepositoryInterface
	ProductDesaServiceInterface            ProductDesaServiceInterface
	OperatorPrefixRepositoryInterface      repository.OperatorPrefixRepositoryInterface
	OrderItemPpobRepositoryInterface       repository.OrderItemPpobRepositoryInterface
	PpobDetailRepositoryInterface          repository.PpobDetailRepositoryInterface
	DesaRepositoryInterface                repository.DesaRepositoryInterface
	InveliAPIRepositoryInterface           invelirepository.InveliAPIRepositoryInterface
	ListPinjamanRepositoryInterface        repository.ListPinjamanRepositoryInterface
	UserShippingAddressRepositoryInterface repository.UserShippingAddressRepositoryInterface
}

func NewOrderService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	orderRepositoryInterface repository.OrderRepositoryInterface,
	userRepositoryInterface repository.UserRepositoryInterface,
	paymentServiceInterface PaymentServiceInterface,
	cartRepositoryInterface repository.CartRepositoryInterface,
	orderItemRepositoryInterface repository.OrderItemRepositoryInterface,
	paymentChannelRepositoryInterface repository.PaymentChannelRepositoryInterface,
	productDesaRepositoryInterface repository.ProductDesaRepositoryInterface,
	productDesaServiceInterface ProductDesaServiceInterface,
	operatorPrefixRepositoryInterface repository.OperatorPrefixRepositoryInterface,
	orderItemPpobRepositoryInterface repository.OrderItemPpobRepositoryInterface,
	ppobDetailRepositoryInterface repository.PpobDetailRepositoryInterface,
	desaRepositoryInterface repository.DesaRepositoryInterface,
	inveliAPIRepositoryInterface invelirepository.InveliAPIRepositoryInterface,
	listPinjamanRepositoryInterface repository.ListPinjamanRepositoryInterface,
	userShippingAddressRepositoryInterface repository.UserShippingAddressRepositoryInterface,
) OrderServiceInterface {
	return &OrderServiceImplementation{
		DB:                                     db,
		Validate:                               validate,
		Logger:                                 logger,
		OrderRepositoryInterface:               orderRepositoryInterface,
		UserRepositoryInterface:                userRepositoryInterface,
		PaymentServiceInterface:                paymentServiceInterface,
		CartRepositoryInterface:                cartRepositoryInterface,
		OrderItemRepositoryInterface:           orderItemRepositoryInterface,
		PaymentChannelRepositoryInterface:      paymentChannelRepositoryInterface,
		ProductDesaRepositoryInterface:         productDesaRepositoryInterface,
		ProductDesaServiceInterface:            productDesaServiceInterface,
		OperatorPrefixRepositoryInterface:      operatorPrefixRepositoryInterface,
		OrderItemPpobRepositoryInterface:       orderItemPpobRepositoryInterface,
		PpobDetailRepositoryInterface:          ppobDetailRepositoryInterface,
		DesaRepositoryInterface:                desaRepositoryInterface,
		InveliAPIRepositoryInterface:           inveliAPIRepositoryInterface,
		ListPinjamanRepositoryInterface:        listPinjamanRepositoryInterface,
		UserShippingAddressRepositoryInterface: userShippingAddressRepositoryInterface,
	}
}

func (service *OrderServiceImplementation) GenerateNumberOrder(idDesa string) (numberOrder string) {
	now := time.Now()
	desa, err := service.DesaRepositoryInterface.FindDesaById(service.DB, idDesa)
	exceptions.PanicIfError(err, "", service.Logger)
	if len(desa.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("desa not found"), "", []string{"desa not found"}, service.Logger)
	}
	for {
		rand.Seed(time.Now().UTC().UnixNano())
		generateCode := 100000 + rand.Intn(999999-100000)
		numberOrder = "ORDER/" + desa.KodeTrx + "/" + now.Format("20060102") + "/" + fmt.Sprint(generateCode)

		// Check number order if exist
		order, _ := service.OrderRepositoryInterface.FindOrderByNumberOrder(service.DB, numberOrder)
		if len(order.Id) == 0 {
			break
		}
	}
	return numberOrder
}

func (service *OrderServiceImplementation) FindOrderPayLaterByIdUser(requestId, idUser string) (orderResponse []response.FindOrderByUserResponse) {
	orders, err := service.OrderRepositoryInterface.FindOrderPayLaterById(service.DB, idUser)
	fmt.Println("masuk service")
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orders) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}
	orderResponse = response.ToFindOrderByUserResponse(orders)
	return orderResponse
}

func (service *OrderServiceImplementation) CreateOrderPostpaidPln(requestId, idUser, idDesa, productType string, orderRequest *request.CreateOrderPostpaidRequest) (createOrderResponse response.CreateOrderResponse) {
	var err error

	request.ValidateRequest(service.Validate, orderRequest, requestId, service.Logger)

	// Get data user
	userProfile, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(userProfile.User.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user not found"), requestId, []string{"user not found"}, service.Logger)
	}

	tx := service.DB.Begin()
	// make object
	orderEntity := &entity.Order{}

	// Generate number and id order
	numberOrder := service.GenerateNumberOrder(idDesa)
	orderEntity.Id = utilities.RandomUUID()
	orderEntity.IdUser = idUser
	orderEntity.IdDesa = idDesa
	orderEntity.RefId = orderRequest.RefId
	orderEntity.NumberOrder = numberOrder
	orderEntity.NamaLengkap = userProfile.NamaLengkap
	orderEntity.Email = userProfile.Email
	orderEntity.Phone = userProfile.User.Phone
	orderEntity.ProductType = productType
	orderEntity.PaymentPoint = orderRequest.PaymentPoint
	orderEntity.OrderedDate = time.Now()
	orderEntity.PaymentMethod = orderRequest.PaymentMethod
	orderEntity.PaymentChannel = orderRequest.PaymentChannel
	orderEntity.TotalBill = orderRequest.TotalBill
	orderEntity.PaymentFee = orderRequest.PaymentFee
	orderEntity.OrderType = 2

	// Check status transaksi
	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + "cs"))
	body, _ := json.Marshal(map[string]interface{}{
		"commands": "checkstatus",
		"ref_id":   orderRequest.RefId,
		"username": config.GetConfig().Ppob.Username,
		"sign":     hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PostpaidUrl
	// fmt.Println("url = ", urlString)
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

	//  cek request

	// reqDump, _ := httputil.DumpRequestOut(req, true)
	// fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	defer resp.Body.Close()

	trxData := &ppob.PostpaidCheckTransactionPln{}
	if err = json.Unmarshal([]byte(data), &trxData); err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	detailTagihan, _ := json.Marshal(trxData.Data.Desc.Tagihan.Detail)

	var totalHarga float64
	var product []string
	var qty []int
	var price []float64

	// data order items ppob
	orderItemsPpob := &entity.OrderItemPpob{}
	orderItemsPpob.Id = utilities.RandomUUID()
	orderItemsPpob.IdOrder = orderEntity.Id
	orderItemsPpob.RefId = orderRequest.RefId
	orderItemsPpob.IdUser = userProfile.IdUser
	orderItemsPpob.ProductCode = trxData.Data.Code
	orderItemsPpob.ProductType = productType
	orderItemsPpob.TotalTagihan = trxData.Data.Price
	orderItemsPpob.Nominal = trxData.Data.Nominal
	orderItemsPpob.Admin = trxData.Data.Admin
	orderItemsPpob.SellingPrice = trxData.Data.SellingPrice
	orderItemsPpob.CreatedAt = time.Now()
	orderItemsPpob.BillDetail = fmt.Sprintf("%s\n", data)

	ppobDetailPln := &entity.PpobDetailPostpaidPln{}
	ppobDetailPln.Id = utilities.RandomUUID()
	ppobDetailPln.IdOrderItemPpob = orderItemsPpob.Id
	ppobDetailPln.RefId = orderRequest.RefId
	ppobDetailPln.CustomerId = orderRequest.CustomerId
	ppobDetailPln.CustomerName = trxData.Data.TrName
	ppobDetailPln.Tarif = trxData.Data.Desc.Tarif
	ppobDetailPln.Daya = trxData.Data.Desc.Daya
	ppobDetailPln.Period = trxData.Data.Period
	ppobDetailPln.LembarTagihan = trxData.Data.Desc.LembarTagihan
	ppobDetailPln.StatusTopUp = -1
	ppobDetailPln.JsonDetailTagihan = fmt.Sprintf("%s\n", detailTagihan)

	// check if cc
	if orderRequest.PaymentMethod == "cc" {
		product = append(product, orderItemsPpob.ProductCode)
		qty = append(qty, 1)
		price = append(price, orderItemsPpob.TotalTagihan)
	}

	totalHarga = trxData.Data.Price

	log.Println("total harga client = ", totalHarga+orderRequest.PaymentFee)
	log.Println("total harga server = ", orderRequest.TotalBill+orderRequest.PaymentFee+orderEntity.PaymentPoint)

	if (totalHarga + orderRequest.PaymentFee) != (orderRequest.TotalBill + orderRequest.PaymentFee + orderEntity.PaymentPoint) {
		exceptions.PanicIfErrorWithRollback(errors.New("harga tidak sama"), requestId, []string{"harga tidak sama"}, service.Logger, tx)
	}

	// Get detail payment channel
	paymentChannel, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(tx, orderRequest.PaymentChannel)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error get payment by code"}, service.Logger, tx)
	if len(paymentChannel.Id) == 0 {
		exceptions.PanicIfRecordNotFoundWithRollback(err, requestId, []string{"payment not found"}, service.Logger, tx)
	}

	// Get Desa
	desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, userProfile.User.IdDesa)
	if len(desa.Id) == 0 {
		exceptions.PanicIfErrorWithRollback(errors.New("desa account paylater not found"), requestId, []string{"desa account paylater not found"}, service.Logger, tx)
	}

	switch orderRequest.PaymentMethod {
	case "point":
		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Point"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
	case "trf":
		// buat nomor acak
		rand.Seed(time.Now().UnixNano())
		min := 111
		max := 299
		rand3Number := rand.Intn(max-min+1) + min

		min2 := 11
		max2 := 99
		rand2Number := rand.Intn(max2-min2+1) + min

		sisaPembagi := math.Mod(orderRequest.TotalBill, 1000)
		var Total float64

		if sisaPembagi < 100 {
			Total = orderRequest.TotalBill + float64(rand3Number)
		} else if sisaPembagi >= 100 {
			Total = orderRequest.TotalBill + float64(rand2Number)
		}

		orderEntity.OrderStatus = 0
		orderEntity.PaymentStatus = 0
		orderEntity.PaymentNo = paymentChannel.NoAccountBank
		orderEntity.PaymentName = paymentChannel.NamaPemilikBank
		orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
		orderEntity.PaymentCash = Total

	case "va", "qris":
		orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		res := service.PaymentServiceInterface.VaQrisPay(requestId,
			&payment.IpaymuQrisVaRequest{
				Name:           userProfile.NamaLengkap,
				Phone:          userProfile.User.Phone,
				Email:          userProfile.Email,
				Amount:         orderEntity.PaymentCash,
				ReferenceId:    numberOrder,
				PaymentMethod:  orderRequest.PaymentMethod,
				PaymentChannel: orderRequest.PaymentChannel,
			},
		)

		if res.Status != 200 {
			fmt.Println("LOG RESPONSE IPAYMU = ", res)
			exceptions.PanicIfErrorWithRollback(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger, tx)
		} else if res.Status == 200 {
			paymentDueDate, _ := time.Parse("2006-01-02 15:04:05", res.Data.Expired)
			orderEntity.PaymentStatus = 0
			orderEntity.TrxId = res.Data.TransactionId
			orderEntity.PaymentNo = res.Data.PaymentNo
			orderEntity.PaymentName = res.Data.PaymentName
			orderEntity.PaymentDueDate = null.NewTime(paymentDueDate, true)
			orderEntity.OrderStatus = 0
		}

	case "cc":
		// tambahkan ongkos kirim
		product = append(product, "Payment Fee")
		qty = append(qty, 1)
		price = append(price, orderRequest.PaymentFee)

		res := service.PaymentServiceInterface.CreditCardPay(requestId,
			&payment.IpaymuCreditCardRequest{
				Product:       product,
				Qty:           qty,
				Price:         price,
				ReferenceId:   numberOrder,
				BuyerName:     userProfile.NamaLengkap,
				BuyerEmail:    userProfile.Email,
				BuyerPhone:    userProfile.User.Phone,
				PaymentMethod: orderRequest.PaymentMethod,
			},
		)

		if res.Status != 200 {
			fmt.Println("LOG RESPONSE IPAYMU = ", res)
			exceptions.PanicIfErrorWithRollback(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger, tx)
		} else if res.Status == 200 {
			orderEntity.PaymentStatus = 0
			orderEntity.PaymentNo = res.Data.Url
			orderEntity.PaymentName = "Credit Card"
			orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
			orderEntity.OrderStatus = 0
			orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		}

	case "paylater":

		var isMerchant float64
		var totalAmount float64

		// Set Is Merchant 0
		isMerchant = 0

		// Validasi Saldo Bupda
		saldoBupda, err := service.InveliAPIRepositoryInterface.GetSaldoBupda(userProfile.User.InveliAccessToken, desa.GroupIdBupda)

		if err != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error saldo bupda "+err.Error()), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		if saldoBupda <= 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("saldo bupda kurang"), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		// Get Bunga
		bunga, errr := service.InveliAPIRepositoryInterface.GetLoanProduct(userProfile.User.InveliAccessToken)
		if errr != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error get loan product "+err.Error()), requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}

		// Get Loan Product
		loandProductID, errr := service.InveliAPIRepositoryInterface.GetLoanProductId(userProfile.User.InveliAccessToken)
		if errr != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error get loan product id "+err.Error()), requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}

		if len(loandProductID) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("loan product id not found"), requestId, []string{"loan product id not found"}, service.Logger, tx)
		}

		// Get Account User
		accountUser, _ := service.UserRepositoryInterface.GetUserAccountPaylaterByID(tx, userProfile.User.Id)
		if len(accountUser.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("user account paylater not found"), requestId, []string{"user account paylater not found"}, service.Logger, tx)
		}

		// Validasi Tunggakan Paylater
		tunggakanPaylater, err := service.InveliAPIRepositoryInterface.GetTunggakan(accountUser.IdAccount, userProfile.User.InveliAccessToken)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(err, requestId, []string{err.Error()}, service.Logger, tx)
		}

		if len(tunggakanPaylater) != 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("masih ada tunggakan"), requestId, []string{"masih ada tunggakan yang belum di bayar"}, service.Logger, tx)
		}

		totalAmount = orderRequest.TotalBill + orderRequest.PaymentFee

		err = service.InveliAPIRepositoryInterface.InveliCreatePaylater(userProfile.User.InveliAccessToken, userProfile.User.InveliIDMember, accountUser.IdAccount, orderRequest.TotalBill, totalAmount, isMerchant, bunga, loandProductID, desa.NoRekening)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error create pinjaman "+err.Error()), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		if time.Now().Local().Day() < 25 {
			orderEntity.PaymentDueDate = null.NewTime(time.Date(time.Now().Year(), time.Now().Month(), 25, 0, 0, 0, 0, time.Local), true)
		} else if time.Now().Local().Day() >= 25 {
			orderEntity.PaymentDueDate = null.NewTime(time.Date(time.Now().Year(), time.Now().Month()+1, 25, 0, 0, 0, 0, time.Local), true)
		}

		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Paylater"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
		orderEntity.PaymentCash = orderRequest.TotalBill + orderRequest.PaymentFee

		// err = service.InveliAPIRepositoryInterface.ApiPayment(desa.NoRekening, accountUser.Code, userProfile.User.InveliAccessToken, orderRequest.TotalBill, isMerchant)
		// if err != nil {
		// 	exceptions.PanicIfErrorWithRollback(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		// }

		var jmlOrder float64
		jmlOrderPayLate, err := service.OrderRepositoryInterface.FindOrderPayLaterById(service.DB, idUser)
		if err != nil {
			log.Println(err.Error())
		}
		jmlOrder = 0
		for _, v := range jmlOrderPayLate {
			jmlOrder = jmlOrder + v.TotalBill
		}

		userPaylaterFlag, _ := service.UserRepositoryInterface.GetUserPayLaterFlagThisMonth(service.DB, idUser)

		if (int(jmlOrder) + int(orderRequest.TotalBill)) > (userPaylaterFlag.TanggungRentengFlag * 1000000) {
			service.UserRepositoryInterface.UpdateUserPayLaterFlag(service.DB, idUser, &entity.UsersPaylaterFlag{
				TanggungRentengFlag: userPaylaterFlag.TanggungRentengFlag + 1,
			})
		}

	case "tabungan_bima":

		accountUser, _ := service.UserRepositoryInterface.GetUserAccountBimaByID(tx, userProfile.User.Id)
		if len(accountUser.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("user account paylater not found"), requestId, []string{"user account paylater not found"}, service.Logger, tx)
		}

		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Tabungan Bima"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
		orderEntity.PaymentCash = orderRequest.TotalBill

		desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, userProfile.User.IdDesa)
		if len(desa.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("desa account paylater not found"), requestId, []string{"desa account paylater not found"}, service.Logger, tx)
		}

		err = service.InveliAPIRepositoryInterface.ApiPayment(desa.NoRekening, accountUser.Code, userProfile.User.InveliAccessToken, orderRequest.TotalBill, 0)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}
	}

	// Create Order
	err = service.OrderRepositoryInterface.CreateOrder(tx, orderEntity)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order"}, service.Logger, tx)

	// Create order items
	err = service.OrderItemPpobRepositoryInterface.CreateOrderItemPpob(tx, orderItemsPpob)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order items"}, service.Logger, tx)

	// create ppob detail pdam
	err = service.PpobDetailRepositoryInterface.CreateOrderPpobDetailPostpaidPln(tx, ppobDetailPln)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order items"}, service.Logger, tx)

	if orderRequest.PaymentMethod == "tabungan_bima" || orderRequest.PaymentMethod == "paylater" {

		response := service.PostpaidTopupPln(requestId, orderRequest.CustomerId, orderItemsPpob.TrId, orderItemsPpob.ProductCode)

		_ = service.PpobDetailRepositoryInterface.UpdatePpobPostpaidPlnById(service.DB, ppobDetailPln.Id, &entity.PpobDetailPostpaidPln{
			StatusTopUp:         3,
			TopupProccesingDate: null.NewTime(time.Now(), true),
			LastBalance:         response.Balance,
		})
	}

	commit := tx.Commit()
	exceptions.PanicIfError(commit.Error, requestId, service.Logger)

	createOrderResponse = response.ToCreateOrderResponse(orderEntity, paymentChannel)
	return createOrderResponse
}

func (service *OrderServiceImplementation) CreateOrderPostpaidPdam(requestId, idUser, idDesa, productType string, orderRequest *request.CreateOrderPostpaidRequest) (createOrderResponse response.CreateOrderResponse) {
	var err error

	request.ValidateRequest(service.Validate, orderRequest, requestId, service.Logger)

	// Get data user
	userProfile, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(userProfile.User.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user not found"), requestId, []string{"user not found"}, service.Logger)
	}

	tx := service.DB.Begin()
	// make object
	orderEntity := &entity.Order{}

	// Generate number and id order
	numberOrder := service.GenerateNumberOrder(idDesa)
	orderEntity.Id = utilities.RandomUUID()
	orderEntity.IdUser = idUser
	orderEntity.IdDesa = idDesa
	orderEntity.RefId = orderRequest.RefId
	orderEntity.NumberOrder = numberOrder
	orderEntity.NamaLengkap = userProfile.NamaLengkap
	orderEntity.Email = userProfile.Email
	orderEntity.Phone = userProfile.User.Phone
	orderEntity.ProductType = productType
	orderEntity.PaymentPoint = orderRequest.PaymentPoint
	orderEntity.OrderedDate = time.Now()
	orderEntity.PaymentMethod = orderRequest.PaymentMethod
	orderEntity.PaymentChannel = orderRequest.PaymentChannel
	orderEntity.TotalBill = orderRequest.TotalBill
	orderEntity.PaymentFee = orderRequest.PaymentFee
	orderEntity.OrderType = 2

	// Check status transaksi
	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + "cs"))
	body, _ := json.Marshal(map[string]interface{}{
		"commands": "checkstatus",
		"ref_id":   orderRequest.RefId,
		"username": config.GetConfig().Ppob.Username,
		"sign":     hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PostpaidUrl
	// fmt.Println("url = ", urlString)
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

	//  cek request

	// reqDump, _ := httputil.DumpRequestOut(req, true)
	// fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	defer resp.Body.Close()

	trxData := &ppob.PostpaidCheckTransactionPdam{}
	if err = json.Unmarshal([]byte(data), &trxData); err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	detailTagihan, _ := json.Marshal(trxData.Data.Desc.Bill.Detail)

	var totalHarga float64
	var product []string
	var qty []int
	var price []float64
	// data order items ppob
	orderItemsPpob := &entity.OrderItemPpob{}
	orderItemsPpob.Id = utilities.RandomUUID()
	orderItemsPpob.IdOrder = orderEntity.Id
	orderItemsPpob.IdUser = userProfile.IdUser
	orderItemsPpob.ProductCode = trxData.Data.Code
	orderItemsPpob.ProductType = productType
	orderItemsPpob.TotalTagihan = trxData.Data.Price
	orderItemsPpob.TrId = trxData.Data.TrxId
	orderItemsPpob.RefId = orderRequest.RefId
	orderItemsPpob.CreatedAt = time.Now()
	orderItemsPpob.BillDetail = fmt.Sprintf("%s\n", data)

	ppobDetailPdam := &entity.PpobDetailPostpaidPdam{}
	ppobDetailPdam.Id = utilities.RandomUUID()
	ppobDetailPdam.IdOrderItemPpob = orderItemsPpob.Id
	ppobDetailPdam.RefId = orderRequest.RefId
	ppobDetailPdam.TrId = orderItemsPpob.TrId
	ppobDetailPdam.CustomerId = trxData.Data.Hp
	ppobDetailPdam.CustomerName = trxData.Data.TrName
	ppobDetailPdam.BillQty = trxData.Data.Desc.BillQuantity
	ppobDetailPdam.Period = trxData.Data.Period
	ppobDetailPdam.DueDate = trxData.Data.Desc.DueDate
	ppobDetailPdam.PdamName = trxData.Data.Desc.PdamName
	ppobDetailPdam.PdamAddress = trxData.Data.Desc.Address
	ppobDetailPdam.StampDuty = trxData.Data.Desc.StampDuty
	ppobDetailPdam.Address = trxData.Data.Desc.Address
	ppobDetailPdam.StatusTopUp = -1
	ppobDetailPdam.JsonDetailTagihan = fmt.Sprintf("%s\n", detailTagihan)

	// check if cc
	if orderRequest.PaymentMethod == "cc" {
		product = append(product, orderItemsPpob.ProductCode)
		qty = append(qty, 1)
		price = append(price, orderItemsPpob.TotalTagihan)
	}

	totalHarga = trxData.Data.Price

	if (totalHarga + orderRequest.PaymentFee) != (orderRequest.TotalBill + orderRequest.PaymentFee + orderEntity.PaymentPoint) {
		exceptions.PanicIfErrorWithRollback(errors.New("harga tidak sama"), requestId, []string{"harga tidak sama"}, service.Logger, tx)
	}

	// Get detail payment channel
	paymentChannel, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(tx, orderRequest.PaymentChannel)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error get payment by code"}, service.Logger, tx)
	if len(paymentChannel.Id) == 0 {
		exceptions.PanicIfRecordNotFoundWithRollback(err, requestId, []string{"payment not found"}, service.Logger, tx)
	}

	// Get Desa
	desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, userProfile.User.IdDesa)
	if len(desa.Id) == 0 {
		exceptions.PanicIfErrorWithRollback(errors.New("desa account paylater not found"), requestId, []string{"desa account paylater not found"}, service.Logger, tx)
	}

	switch orderRequest.PaymentMethod {
	case "point":
		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Point"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
	case "trf":
		// buat nomor acak
		rand.Seed(time.Now().UnixNano())
		min := 111
		max := 299
		rand3Number := rand.Intn(max-min+1) + min

		min2 := 11
		max2 := 99
		rand2Number := rand.Intn(max2-min2+1) + min

		sisaPembagi := math.Mod(orderRequest.TotalBill, 1000)
		var Total float64

		if sisaPembagi < 100 {
			Total = orderRequest.TotalBill + float64(rand3Number)
		} else if sisaPembagi >= 100 {
			Total = orderRequest.TotalBill + float64(rand2Number)
		}

		orderEntity.OrderStatus = 0
		orderEntity.PaymentStatus = 0
		orderEntity.PaymentNo = paymentChannel.NoAccountBank
		orderEntity.PaymentName = paymentChannel.NamaPemilikBank
		orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
		orderEntity.PaymentCash = Total

	case "va", "qris":
		orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		res := service.PaymentServiceInterface.VaQrisPay(requestId,
			&payment.IpaymuQrisVaRequest{
				Name:           userProfile.NamaLengkap,
				Phone:          userProfile.User.Phone,
				Email:          userProfile.Email,
				Amount:         orderEntity.PaymentCash,
				ReferenceId:    numberOrder,
				PaymentMethod:  orderRequest.PaymentMethod,
				PaymentChannel: orderRequest.PaymentChannel,
			},
		)

		if res.Status != 200 {
			fmt.Println("LOG RESPONSE IPAYMU = ", res)
			exceptions.PanicIfErrorWithRollback(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger, tx)
		} else if res.Status == 200 {
			paymentDueDate, _ := time.Parse("2006-01-02 15:04:05", res.Data.Expired)
			orderEntity.PaymentStatus = 0
			orderEntity.TrxId = res.Data.TransactionId
			orderEntity.PaymentNo = res.Data.PaymentNo
			orderEntity.PaymentName = res.Data.PaymentName
			orderEntity.PaymentDueDate = null.NewTime(paymentDueDate, true)
			orderEntity.OrderStatus = 0
		}

	case "cc":
		// tambahkan ongkos kirim
		product = append(product, "Payment Fee")
		qty = append(qty, 1)
		price = append(price, orderRequest.PaymentFee)

		res := service.PaymentServiceInterface.CreditCardPay(requestId,
			&payment.IpaymuCreditCardRequest{
				Product:       product,
				Qty:           qty,
				Price:         price,
				ReferenceId:   numberOrder,
				BuyerName:     userProfile.NamaLengkap,
				BuyerEmail:    userProfile.Email,
				BuyerPhone:    userProfile.User.Phone,
				PaymentMethod: orderRequest.PaymentMethod,
			},
		)

		if res.Status != 200 {
			fmt.Println("LOG RESPONSE IPAYMU = ", res)
			exceptions.PanicIfErrorWithRollback(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger, tx)
		} else if res.Status == 200 {
			orderEntity.PaymentStatus = 0
			orderEntity.PaymentNo = res.Data.Url
			orderEntity.PaymentName = "Credit Card"
			orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
			orderEntity.OrderStatus = 0
			orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		}

	case "paylater":
		var isMerchant float64
		var totalAmount float64

		// Set Is Merchant 0
		isMerchant = 0

		// Validasi Saldo Bupda
		saldoBupda, err := service.InveliAPIRepositoryInterface.GetSaldoBupda(userProfile.User.InveliAccessToken, desa.GroupIdBupda)

		if err != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error saldo bupda "+err.Error()), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		if saldoBupda <= 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("saldo bupda kurang"), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		// Get Bunga
		bunga, errr := service.InveliAPIRepositoryInterface.GetLoanProduct(userProfile.User.InveliAccessToken)
		if errr != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error get loan product "+err.Error()), requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}

		// Get Loan Product
		loandProductID, errr := service.InveliAPIRepositoryInterface.GetLoanProductId(userProfile.User.InveliAccessToken)
		if errr != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error get loan product id "+err.Error()), requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}

		if len(loandProductID) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("loan product id not found"), requestId, []string{"loan product id not found"}, service.Logger, tx)
		}

		// Get Account User
		accountUser, _ := service.UserRepositoryInterface.GetUserAccountPaylaterByID(tx, userProfile.User.Id)
		if len(accountUser.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("user account paylater not found"), requestId, []string{"user account paylater not found"}, service.Logger, tx)
		}

		// Validasi Tunggakan Paylater
		tunggakanPaylater, err := service.InveliAPIRepositoryInterface.GetTunggakan(accountUser.IdAccount, userProfile.User.InveliAccessToken)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(err, requestId, []string{err.Error()}, service.Logger, tx)
		}

		if len(tunggakanPaylater) != 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("masih ada tunggakan"), requestId, []string{"masih ada tunggakan yang belum di bayar"}, service.Logger, tx)
		}

		totalAmount = orderRequest.TotalBill + orderRequest.PaymentFee

		err = service.InveliAPIRepositoryInterface.InveliCreatePaylater(userProfile.User.InveliAccessToken, userProfile.User.InveliIDMember, accountUser.IdAccount, orderRequest.TotalBill, totalAmount, isMerchant, bunga, loandProductID, desa.NoRekening)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error care pinjaman "+err.Error()), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		if time.Now().Local().Day() < 25 {
			orderEntity.PaymentDueDate = null.NewTime(time.Date(time.Now().Year(), time.Now().Month(), 25, 0, 0, 0, 0, time.Local), true)
		} else if time.Now().Local().Day() >= 25 {
			orderEntity.PaymentDueDate = null.NewTime(time.Date(time.Now().Year(), time.Now().Month()+1, 25, 0, 0, 0, 0, time.Local), true)
		}

		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Paylater"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
		orderEntity.PaymentCash = orderRequest.TotalBill + orderRequest.PaymentFee

		// err = service.InveliAPIRepositoryInterface.ApiPayment(desa.NoRekening, accountUser.Code, userProfile.User.InveliAccessToken, orderRequest.TotalBill, isMerchant)
		// if err != nil {
		// 	exceptions.PanicIfErrorWithRollback(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		// }

		var jmlOrder float64
		jmlOrderPayLate, err := service.OrderRepositoryInterface.FindOrderPayLaterById(service.DB, idUser)
		if err != nil {
			log.Println(err.Error())
		}
		jmlOrder = 0
		for _, v := range jmlOrderPayLate {
			jmlOrder = jmlOrder + v.TotalBill
		}

		userPaylaterFlag, _ := service.UserRepositoryInterface.GetUserPayLaterFlagThisMonth(service.DB, idUser)

		if (int(jmlOrder) + int(orderRequest.TotalBill)) > (userPaylaterFlag.TanggungRentengFlag * 1000000) {
			service.UserRepositoryInterface.UpdateUserPayLaterFlag(service.DB, idUser, &entity.UsersPaylaterFlag{
				TanggungRentengFlag: userPaylaterFlag.TanggungRentengFlag + 1,
			})
		}

	case "tabungan_bima":

		accountUser, _ := service.UserRepositoryInterface.GetUserAccountBimaByID(tx, userProfile.User.Id)
		if len(accountUser.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("user account paylater not found"), requestId, []string{"user account paylater not found"}, service.Logger, tx)
		}

		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Tabungan Bima"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
		orderEntity.PaymentCash = orderRequest.TotalBill

		desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, userProfile.User.IdDesa)
		if len(desa.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("desa account paylater not found"), requestId, []string{"desa account paylater not found"}, service.Logger, tx)
		}

		err = service.InveliAPIRepositoryInterface.ApiPayment(desa.NoRekening, accountUser.Code, userProfile.User.InveliAccessToken, orderRequest.TotalBill, 0)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}
	}

	// Create Order
	err = service.OrderRepositoryInterface.CreateOrder(tx, orderEntity)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order"}, service.Logger, tx)

	// Create order items
	err = service.OrderItemPpobRepositoryInterface.CreateOrderItemPpob(tx, orderItemsPpob)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order items"}, service.Logger, tx)

	// create ppob detail pdam
	err = service.PpobDetailRepositoryInterface.CreateOrderPpobDetailPostpaidPdam(tx, ppobDetailPdam)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order items"}, service.Logger, tx)

	if orderRequest.PaymentMethod == "tabungan_bima" || orderRequest.PaymentMethod == "paylater" {
		response := service.PostpaidTopupPdam(requestId, orderRequest.CustomerId, orderItemsPpob.TrId, orderItemsPpob.ProductCode)

		_ = service.PpobDetailRepositoryInterface.UpdatePpobPostpaidPdamById(service.DB, ppobDetailPdam.Id, &entity.PpobDetailPostpaidPdam{
			StatusTopUp:         3,
			TopupProccesingDate: null.NewTime(time.Now(), true),
			LastBalance:         response.Balance,
		})
	}

	commit := tx.Commit()
	exceptions.PanicIfError(commit.Error, requestId, service.Logger)

	createOrderResponse = response.ToCreateOrderResponse(orderEntity, paymentChannel)
	return createOrderResponse
}

func (service *OrderServiceImplementation) CreateOrderPostpaidTelco(requestId, idUser, idDesa, productType string, orderRequest *request.CreateOrderPostpaidRequest) (createOrderResponse response.CreateOrderResponse) {
	var err error

	request.ValidateRequest(service.Validate, orderRequest, requestId, service.Logger)

	// Get data user
	userProfile, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(userProfile.User.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user not found"), requestId, []string{"user not found"}, service.Logger)
	}

	tx := service.DB.Begin()
	// make object
	orderEntity := &entity.Order{}

	// Generate number and id order
	numberOrder := service.GenerateNumberOrder(idDesa)
	orderEntity.Id = utilities.RandomUUID()
	orderEntity.IdUser = idUser
	orderEntity.IdDesa = idDesa
	orderEntity.RefId = orderRequest.RefId
	orderEntity.NumberOrder = numberOrder
	orderEntity.NamaLengkap = userProfile.NamaLengkap
	orderEntity.Email = userProfile.Email
	orderEntity.Phone = userProfile.User.Phone
	orderEntity.ProductType = productType
	orderEntity.PaymentPoint = orderRequest.PaymentPoint
	orderEntity.OrderedDate = time.Now()
	orderEntity.PaymentMethod = orderRequest.PaymentMethod
	orderEntity.PaymentChannel = orderRequest.PaymentChannel
	orderEntity.TotalBill = orderRequest.TotalBill
	orderEntity.PaymentFee = orderRequest.PaymentFee
	orderEntity.OrderType = 2

	// Check status transaksi
	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + "cs"))
	body, _ := json.Marshal(map[string]interface{}{
		"commands": "checkstatus",
		"ref_id":   orderRequest.RefId,
		"username": config.GetConfig().Ppob.Username,
		"sign":     hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PostpaidUrl
	// fmt.Println("url = ", urlString)
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

	// cek request

	// reqDump, _ := httputil.DumpRequestOut(req, true)
	// fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	defer resp.Body.Close()

	trxData := &ppob.PostpaidCheckTransactionTelco{}
	if err = json.Unmarshal([]byte(data), &trxData); err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	detailTagihan, _ := json.Marshal(trxData.Data.Desc.Tagihan)

	var totalHarga float64
	var product []string
	var qty []int
	var price []float64
	// data order items ppob
	orderItemsPpob := &entity.OrderItemPpob{}
	orderItemsPpob.Id = utilities.RandomUUID()
	orderItemsPpob.IdOrder = orderEntity.Id
	orderItemsPpob.IdUser = userProfile.IdUser
	orderItemsPpob.ProductCode = trxData.Data.Code
	orderItemsPpob.ProductType = productType
	orderItemsPpob.TotalTagihan = trxData.Data.Price
	orderItemsPpob.TrId = trxData.Data.TrxId
	orderItemsPpob.RefId = orderRequest.RefId
	orderItemsPpob.CreatedAt = time.Now()
	orderItemsPpob.BillDetail = fmt.Sprintf("%s\n", data)

	ppobDetailTelco := &entity.PpobDetailPostpaidTelco{}
	ppobDetailTelco.Id = utilities.RandomUUID()
	ppobDetailTelco.IdOrderItemPpob = orderItemsPpob.Id
	ppobDetailTelco.RefId = orderRequest.RefId
	ppobDetailTelco.TrId = orderItemsPpob.TrId
	ppobDetailTelco.CustomerId = trxData.Data.Hp
	ppobDetailTelco.CustomerName = trxData.Data.TrName
	ppobDetailTelco.Period = trxData.Data.Period
	ppobDetailTelco.KodeArea = trxData.Data.Desc.KodeArea
	ppobDetailTelco.Divre = trxData.Data.Desc.Divre
	ppobDetailTelco.Datel = trxData.Data.Desc.Datel
	ppobDetailTelco.StatusTopUp = -1
	ppobDetailTelco.JsonDetailTagihan = fmt.Sprintf("%s\n", detailTagihan)

	// check if cc
	if orderRequest.PaymentMethod == "cc" {
		product = append(product, orderItemsPpob.ProductCode)
		qty = append(qty, 1)
		price = append(price, orderItemsPpob.TotalTagihan)
	}

	totalHarga = trxData.Data.Price

	if (totalHarga + orderRequest.PaymentFee) != (orderRequest.TotalBill + orderRequest.PaymentFee + orderEntity.PaymentPoint) {
		exceptions.PanicIfErrorWithRollback(errors.New("harga tidak sama"), requestId, []string{"harga tidak sama"}, service.Logger, tx)
	}

	// Get detail payment channel
	paymentChannel, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(tx, orderRequest.PaymentChannel)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error get payment by code"}, service.Logger, tx)
	if len(paymentChannel.Id) == 0 {
		exceptions.PanicIfRecordNotFoundWithRollback(err, requestId, []string{"payment not found"}, service.Logger, tx)
	}

	// Get Desa
	desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, userProfile.User.IdDesa)
	if len(desa.Id) == 0 {
		exceptions.PanicIfErrorWithRollback(errors.New("desa account paylater not found"), requestId, []string{"desa account paylater not found"}, service.Logger, tx)
	}

	switch orderRequest.PaymentMethod {
	case "point":
		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Point"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
	case "trf":
		// buat nomor acak
		rand.Seed(time.Now().UnixNano())
		min := 111
		max := 299
		rand3Number := rand.Intn(max-min+1) + min

		min2 := 11
		max2 := 99
		rand2Number := rand.Intn(max2-min2+1) + min

		sisaPembagi := math.Mod(orderRequest.TotalBill, 1000)
		var Total float64

		if sisaPembagi < 100 {
			Total = orderRequest.TotalBill + float64(rand3Number)
		} else if sisaPembagi >= 100 {
			Total = orderRequest.TotalBill + float64(rand2Number)
		}

		orderEntity.OrderStatus = 0
		orderEntity.PaymentStatus = 0
		orderEntity.PaymentNo = paymentChannel.NoAccountBank
		orderEntity.PaymentName = paymentChannel.NamaPemilikBank
		orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
		orderEntity.PaymentCash = Total

	case "va", "qris":
		orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		res := service.PaymentServiceInterface.VaQrisPay(requestId,
			&payment.IpaymuQrisVaRequest{
				Name:           userProfile.NamaLengkap,
				Phone:          userProfile.User.Phone,
				Email:          userProfile.Email,
				Amount:         orderEntity.PaymentCash,
				ReferenceId:    numberOrder,
				PaymentMethod:  orderRequest.PaymentMethod,
				PaymentChannel: orderRequest.PaymentChannel,
			},
		)

		if res.Status != 200 {
			fmt.Println("LOG RESPONSE IPAYMU = ", res)
			exceptions.PanicIfErrorWithRollback(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger, tx)
		} else if res.Status == 200 {
			paymentDueDate, _ := time.Parse("2006-01-02 15:04:05", res.Data.Expired)
			orderEntity.PaymentStatus = 0
			orderEntity.TrxId = res.Data.TransactionId
			orderEntity.PaymentNo = res.Data.PaymentNo
			orderEntity.PaymentName = res.Data.PaymentName
			orderEntity.PaymentDueDate = null.NewTime(paymentDueDate, true)
			orderEntity.OrderStatus = 0
		}

	case "cc":
		// tambahkan ongkos kirim
		product = append(product, "Payment Fee")
		qty = append(qty, 1)
		price = append(price, orderRequest.PaymentFee)

		res := service.PaymentServiceInterface.CreditCardPay(requestId,
			&payment.IpaymuCreditCardRequest{
				Product:       product,
				Qty:           qty,
				Price:         price,
				ReferenceId:   numberOrder,
				BuyerName:     userProfile.NamaLengkap,
				BuyerEmail:    userProfile.Email,
				BuyerPhone:    userProfile.User.Phone,
				PaymentMethod: orderRequest.PaymentMethod,
			},
		)

		if res.Status != 200 {
			fmt.Println("LOG RESPONSE IPAYMU = ", res)
			exceptions.PanicIfErrorWithRollback(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger, tx)
		} else if res.Status == 200 {
			orderEntity.PaymentStatus = 0
			orderEntity.PaymentNo = res.Data.Url
			orderEntity.PaymentName = "Credit Card"
			orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
			orderEntity.OrderStatus = 0
			orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		}

	case "paylater":
		var isMerchant float64
		var totalAmount float64

		// Set Is Merchant 0
		isMerchant = 0

		// Validasi Saldo Bupda
		saldoBupda, err := service.InveliAPIRepositoryInterface.GetSaldoBupda(userProfile.User.InveliAccessToken, desa.GroupIdBupda)

		if err != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error saldo bupda "+err.Error()), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		if saldoBupda <= 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("saldo bupda kurang"), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		// Get Bunga
		bunga, errr := service.InveliAPIRepositoryInterface.GetLoanProduct(userProfile.User.InveliAccessToken)
		if errr != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error get loan product "+err.Error()), requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}

		// Get Loan Product
		loandProductID, errr := service.InveliAPIRepositoryInterface.GetLoanProductId(userProfile.User.InveliAccessToken)
		if errr != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error get loan product id "+err.Error()), requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}

		if len(loandProductID) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("loan product id not found"), requestId, []string{"loan product id not found"}, service.Logger, tx)
		}

		// Get Account User
		accountUser, _ := service.UserRepositoryInterface.GetUserAccountPaylaterByID(tx, userProfile.User.Id)
		if len(accountUser.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("user account paylater not found"), requestId, []string{"user account paylater not found"}, service.Logger, tx)
		}

		// Validasi Tunggakan Paylater
		tunggakanPaylater, err := service.InveliAPIRepositoryInterface.GetTunggakan(accountUser.IdAccount, userProfile.User.InveliAccessToken)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(err, requestId, []string{err.Error()}, service.Logger, tx)
		}

		if len(tunggakanPaylater) != 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("masih ada tunggakan"), requestId, []string{"masih ada tunggakan yang belum di bayar"}, service.Logger, tx)
		}

		totalAmount = orderRequest.TotalBill + orderRequest.PaymentFee

		err = service.InveliAPIRepositoryInterface.InveliCreatePaylater(userProfile.User.InveliAccessToken, userProfile.User.InveliIDMember, accountUser.IdAccount, orderRequest.TotalBill, totalAmount, isMerchant, bunga, loandProductID, desa.NoRekening)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error care pinjaman "+err.Error()), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		if time.Now().Local().Day() < 25 {
			orderEntity.PaymentDueDate = null.NewTime(time.Date(time.Now().Year(), time.Now().Month(), 25, 0, 0, 0, 0, time.Local), true)
		} else if time.Now().Local().Day() >= 25 {
			orderEntity.PaymentDueDate = null.NewTime(time.Date(time.Now().Year(), time.Now().Month()+1, 25, 0, 0, 0, 0, time.Local), true)
		}

		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Paylater"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
		orderEntity.PaymentCash = orderRequest.TotalBill + orderRequest.PaymentFee

		// err = service.InveliAPIRepositoryInterface.ApiPayment(desa.NoRekening, accountUser.Code, userProfile.User.InveliAccessToken, orderRequest.TotalBill, isMerchant)
		// if err != nil {
		// 	exceptions.PanicIfErrorWithRollback(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		// }

		var jmlOrder float64
		jmlOrderPayLate, err := service.OrderRepositoryInterface.FindOrderPayLaterById(service.DB, idUser)
		if err != nil {
			log.Println(err.Error())
		}
		jmlOrder = 0
		for _, v := range jmlOrderPayLate {
			jmlOrder = jmlOrder + v.TotalBill
		}

		userPaylaterFlag, _ := service.UserRepositoryInterface.GetUserPayLaterFlagThisMonth(service.DB, idUser)

		if (int(jmlOrder) + int(orderRequest.TotalBill)) > (userPaylaterFlag.TanggungRentengFlag * 1000000) {
			service.UserRepositoryInterface.UpdateUserPayLaterFlag(service.DB, idUser, &entity.UsersPaylaterFlag{
				TanggungRentengFlag: userPaylaterFlag.TanggungRentengFlag + 1,
			})
		}

	case "tabungan_bima":

		accountUser, _ := service.UserRepositoryInterface.GetUserAccountBimaByID(tx, userProfile.User.Id)
		if len(accountUser.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("user account paylater not found"), requestId, []string{"user account paylater not found"}, service.Logger, tx)
		}

		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Tabungan Bima"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
		orderEntity.PaymentCash = orderRequest.TotalBill

		desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, userProfile.User.IdDesa)
		if len(desa.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("desa account paylater not found"), requestId, []string{"desa account paylater not found"}, service.Logger, tx)
		}

		err = service.InveliAPIRepositoryInterface.ApiPayment(desa.NoRekening, accountUser.Code, userProfile.User.InveliAccessToken, orderRequest.TotalBill, 0)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}
	}

	// Create Order
	err = service.OrderRepositoryInterface.CreateOrder(tx, orderEntity)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order"}, service.Logger, tx)

	// Create order items
	err = service.OrderItemPpobRepositoryInterface.CreateOrderItemPpob(tx, orderItemsPpob)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order items"}, service.Logger, tx)

	// create ppob detail pdam
	err = service.PpobDetailRepositoryInterface.CreateOrderPpobDetailPostpaidTelco(tx, ppobDetailTelco)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order items"}, service.Logger, tx)

	commit := tx.Commit()
	exceptions.PanicIfError(commit.Error, requestId, service.Logger)

	createOrderResponse = response.ToCreateOrderResponse(orderEntity, paymentChannel)
	return createOrderResponse
}

func (service *OrderServiceImplementation) CreateOrderPrepaidPulsa(requestId, idUser, idDesa, productType string, orderRequest *request.CreateOrderPrepaidRequest) (createOrderResponse response.CreateOrderResponse) {

	var err error

	request.ValidateRequest(service.Validate, orderRequest, requestId, service.Logger)

	// Get data user
	userProfile, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(userProfile.User.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user not found"), requestId, []string{"user not found"}, service.Logger)
	}

	tx := service.DB.Begin()
	// make object
	orderEntity := &entity.Order{}

	// Generate number and id order
	numberOrder := service.GenerateNumberOrder(idDesa)
	orderEntity.Id = utilities.RandomUUID()
	orderEntity.IdUser = idUser
	orderEntity.IdDesa = idDesa
	orderEntity.NumberOrder = numberOrder
	orderEntity.NamaLengkap = userProfile.NamaLengkap
	orderEntity.Email = userProfile.Email
	orderEntity.Phone = userProfile.User.Phone
	orderEntity.ProductType = productType
	orderEntity.PaymentPoint = orderRequest.PaymentPoint
	orderEntity.OrderedDate = time.Now()
	orderEntity.PaymentMethod = orderRequest.PaymentMethod
	orderEntity.PaymentChannel = orderRequest.PaymentChannel
	orderEntity.TotalBill = orderRequest.TotalBill
	orderEntity.PaymentFee = orderRequest.PaymentFee
	orderEntity.OrderType = 2
	orderEntity.RefId = utilities.GenerateRefId()

	var operator string

	phone := PrefixNumber(orderRequest.CustomerId)
	opereratorPrefixResult, err := service.OperatorPrefixRepositoryInterface.FindOperatorPrefixByPhone(service.DB, phone)
	exceptions.PanicIfError(err, requestId, service.Logger)
	operator = opereratorPrefixResult.KodeOperator

	// Get Data from iak
	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + "pl"))
	body, _ := json.Marshal(map[string]interface{}{
		"status":   "all",
		"username": config.GetConfig().Ppob.Username,
		"sign":     hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	var typePpob string
	if productType == "prepaid_pulsa" {
		typePpob = "pulsa"
	} else if productType == "prepaid_data" {
		typePpob = "data"
	}

	urlString := config.GetConfig().Ppob.PrepaidHost + "/pricelist/" + typePpob + "/" + operator
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

	//  cek request

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	// fmt.Printf("body: %s\n", data)

	defer resp.Body.Close()

	priceLists := &ppob.PrepaidPriceListResponse{}
	if err = json.Unmarshal([]byte(data), &priceLists); err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	var totalHarga float64
	var product []string
	var qty []int
	var price []float64
	orderItemsPpob := &entity.OrderItemPpob{}
	ppobDetailPrepaidPulsa := &entity.PpobDetailPrepaidPulsa{}
	for _, priceList := range priceLists.Data.Data {
		if priceList.ProductCode == orderRequest.ProductCode {
			totalHarga = priceList.ProductPrice
			orderItemsPpob.Id = utilities.RandomUUID()
			orderItemsPpob.IdOrder = orderEntity.Id
			orderItemsPpob.IdUser = userProfile.IdUser
			orderItemsPpob.RefId = orderEntity.RefId
			orderItemsPpob.ProductCode = priceList.ProductCode
			orderItemsPpob.ProductType = productType
			orderItemsPpob.TotalTagihan = priceList.ProductPrice
			orderItemsPpob.IconUrl = priceList.IconUrl
			orderItemsPpob.CreatedAt = time.Now()
			orderItemsPpob.BillDetail = fmt.Sprintf("%s\n", data)

			ppobDetailPrepaidPulsa.Id = utilities.RandomUUID()
			ppobDetailPrepaidPulsa.IdOrderItemPpob = orderItemsPpob.Id
			ppobDetailPrepaidPulsa.ProductCode = priceList.ProductCode
			ppobDetailPrepaidPulsa.ProductName = priceList.ProductNominal
			ppobDetailPrepaidPulsa.ProductDescription = priceList.ProductDescription
			ppobDetailPrepaidPulsa.CustomerId = orderRequest.CustomerId
			ppobDetailPrepaidPulsa.Operator = operator
			ppobDetailPrepaidPulsa.ActivePeriod = priceList.ActivePeriod
			ppobDetailPrepaidPulsa.IconUrl = priceList.IconUrl
			ppobDetailPrepaidPulsa.StatusTopUp = -1

			if orderRequest.PaymentMethod == "cc" {
				product = append(product, orderItemsPpob.ProductCode)
				qty = append(qty, 1)
				price = append(price, orderItemsPpob.TotalTagihan)
			}
			break
		}
	}

	fmt.Println("Total Harga = ", totalHarga)
	fmt.Println("Total Harga request = ", orderRequest.TotalBill)

	if ((totalHarga + 1500) + orderRequest.PaymentFee) != (orderRequest.TotalBill + orderRequest.PaymentFee + orderEntity.PaymentPoint) {
		exceptions.PanicIfErrorWithRollback(errors.New("harga tidak sama"), requestId, []string{"harga tidak sama"}, service.Logger, tx)
	}

	// Get detail payment channel
	paymentChannel, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(tx, orderRequest.PaymentChannel)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error get payment by code"}, service.Logger, tx)
	if len(paymentChannel.Id) == 0 {
		exceptions.PanicIfRecordNotFoundWithRollback(err, requestId, []string{"payment not found"}, service.Logger, tx)
	}

	// Get Desa
	desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, userProfile.User.IdDesa)
	if len(desa.Id) == 0 {
		exceptions.PanicIfErrorWithRollback(errors.New("desa account paylater not found"), requestId, []string{"desa account paylater not found"}, service.Logger, tx)
	}

	switch orderRequest.PaymentMethod {
	case "point":
		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Point"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
	case "trf":
		// buat nomor acak
		rand.Seed(time.Now().UnixNano())
		min := 111
		max := 299
		rand3Number := rand.Intn(max-min+1) + min

		min2 := 11
		max2 := 99
		rand2Number := rand.Intn(max2-min2+1) + min

		sisaPembagi := math.Mod(orderRequest.TotalBill, 1000)
		var Total float64

		if sisaPembagi < 100 {
			Total = orderRequest.TotalBill + float64(rand3Number)
		} else if sisaPembagi >= 100 {
			Total = orderRequest.TotalBill + float64(rand2Number)
		}

		orderEntity.OrderStatus = 0
		orderEntity.PaymentStatus = 0
		orderEntity.PaymentNo = paymentChannel.NoAccountBank
		orderEntity.PaymentName = paymentChannel.NamaPemilikBank
		orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
		orderEntity.PaymentCash = Total

	case "va", "qris":
		orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		res := service.PaymentServiceInterface.VaQrisPay(requestId,
			&payment.IpaymuQrisVaRequest{
				Name:           userProfile.NamaLengkap,
				Phone:          userProfile.User.Phone,
				Email:          userProfile.Email,
				Amount:         orderEntity.PaymentCash,
				ReferenceId:    numberOrder,
				PaymentMethod:  orderRequest.PaymentMethod,
				PaymentChannel: orderRequest.PaymentChannel,
			},
		)

		if res.Status != 200 {
			fmt.Println("LOG RESPONSE IPAYMU = ", res)
			exceptions.PanicIfErrorWithRollback(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger, tx)
		} else if res.Status == 200 {
			paymentDueDate, _ := time.Parse("2006-01-02 15:04:05", res.Data.Expired)
			orderEntity.PaymentStatus = 0
			orderEntity.TrxId = res.Data.TransactionId
			orderEntity.PaymentNo = res.Data.PaymentNo
			orderEntity.PaymentName = res.Data.PaymentName
			orderEntity.PaymentDueDate = null.NewTime(paymentDueDate, true)
			orderEntity.OrderStatus = 0
		}

	case "cc":
		// tambahkan ongkos kirim
		product = append(product, "Payment Fee")
		qty = append(qty, 1)
		price = append(price, orderRequest.PaymentFee)

		res := service.PaymentServiceInterface.CreditCardPay(requestId,
			&payment.IpaymuCreditCardRequest{
				Product:       product,
				Qty:           qty,
				Price:         price,
				ReferenceId:   numberOrder,
				BuyerName:     userProfile.NamaLengkap,
				BuyerEmail:    userProfile.Email,
				BuyerPhone:    userProfile.User.Phone,
				PaymentMethod: orderRequest.PaymentMethod,
			},
		)

		if res.Status != 200 {
			fmt.Println("LOG RESPONSE IPAYMU = ", res)
			exceptions.PanicIfErrorWithRollback(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger, tx)
		} else if res.Status == 200 {
			orderEntity.PaymentStatus = 0
			orderEntity.PaymentNo = res.Data.Url
			orderEntity.PaymentName = "Credit Card"
			orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
			orderEntity.OrderStatus = 0
			orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		}

	case "paylater":
		var isMerchant float64
		var totalAmount float64

		// Set Is Merchant 0
		isMerchant = 0

		// Validasi Saldo Bupda
		saldoBupda, err := service.InveliAPIRepositoryInterface.GetSaldoBupda(userProfile.User.InveliAccessToken, desa.GroupIdBupda)

		if err != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error saldo bupda "+err.Error()), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		if saldoBupda <= 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("saldo bupda kurang"), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		if saldoBupda <= (orderRequest.TotalBill + orderRequest.PaymentFee) {
			exceptions.PanicIfErrorWithRollback(errors.New("saldo bupda kurang"), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		// Get Bunga
		bunga, errr := service.InveliAPIRepositoryInterface.GetLoanProduct(userProfile.User.InveliAccessToken)
		if errr != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error get loan product "+err.Error()), requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}

		// Get Loan Product
		loandProductID, errr := service.InveliAPIRepositoryInterface.GetLoanProductId(userProfile.User.InveliAccessToken)
		if errr != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error get loan product id "+err.Error()), requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}

		if len(loandProductID) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("loan product id not found"), requestId, []string{"loan product id not found"}, service.Logger, tx)
		}

		// Get Account User
		accountUser, _ := service.UserRepositoryInterface.GetUserAccountPaylaterByID(tx, userProfile.User.Id)
		if len(accountUser.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("user account paylater not found"), requestId, []string{"user account paylater not found"}, service.Logger, tx)
		}

		// Validasi Tunggakan Paylater
		// tunggakanPaylater, err := service.InveliAPIRepositoryInterface.GetTunggakan(accountUser.IdAccount, userProfile.User.InveliAccessToken)
		// if err != nil {
		// 	exceptions.PanicIfErrorWithRollback(err, requestId, []string{err.Error()}, service.Logger, tx)
		// }

		// if len(tunggakanPaylater) != 0 {
		// 	exceptions.PanicIfErrorWithRollback(errors.New("masih ada tunggakan"), requestId, []string{"masih ada tunggakan yang belum di bayar"}, service.Logger, tx)
		// }

		loanID, err := service.InveliAPIRepositoryInterface.GetRiwayatPinjaman(userProfile.User.InveliAccessToken, userProfile.User.InveliIDMember)
		if err != nil {
			log.Println("error get riwayat pinjaman", err.Error())
			exceptions.PanicIfError(err, requestId, service.Logger)
		}

		if len(loanID) != 0 {
			exceptions.PanicIfBadRequest(errors.New("masih ada tunggakan"), requestId, []string{"anda masih memiliki tunggakan"}, service.Logger)
		}

		totalAmount = orderRequest.TotalBill + orderRequest.PaymentFee

		err = service.InveliAPIRepositoryInterface.InveliCreatePaylater(userProfile.User.InveliAccessToken, userProfile.User.InveliIDMember, accountUser.IdAccount, orderRequest.TotalBill, totalAmount, isMerchant, bunga, loandProductID, desa.NoRekening)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error care pinjaman "+err.Error()), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		// if time.Now().Local().Day() < 25 {
		// 	orderEntity.PaymentDueDate = null.NewTime(time.Date(time.Now().Year(), time.Now().Month(), 25, 0, 0, 0, 0, time.Local), true)
		// } else if time.Now().Local().Day() >= 25 {
		// 	orderEntity.PaymentDueDate = null.NewTime(time.Date(time.Now().Year(), time.Now().Month()+1, 25, 0, 0, 0, 0, time.Local), true)
		// }

		orderEntity.PaymentDueDate = null.NewTime(time.Now().AddDate(0, 0, 30), true)

		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Paylater"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
		orderEntity.PaymentCash = orderRequest.TotalBill + orderRequest.PaymentFee

		// err = service.InveliAPIRepositoryInterface.ApiPayment(desa.NoRekening, accountUser.Code, userProfile.User.InveliAccessToken, orderRequest.TotalBill, isMerchant)
		// if err != nil {
		// 	exceptions.PanicIfErrorWithRollback(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		// }

		var jmlOrder float64
		jmlOrderPayLate, err := service.OrderRepositoryInterface.FindOrderPayLaterById(service.DB, idUser)
		if err != nil {
			log.Println(err.Error())
		}
		jmlOrder = 0
		for _, v := range jmlOrderPayLate {
			jmlOrder = jmlOrder + v.TotalBill
		}

		userPaylaterFlag, _ := service.UserRepositoryInterface.GetUserPayLaterFlagThisMonth(service.DB, idUser)

		if (int(jmlOrder) + int(orderRequest.TotalBill)) > (userPaylaterFlag.TanggungRentengFlag * 1000000) {
			service.UserRepositoryInterface.UpdateUserPayLaterFlag(service.DB, idUser, &entity.UsersPaylaterFlag{
				TanggungRentengFlag: userPaylaterFlag.TanggungRentengFlag + 1,
			})
		}

	case "tabungan_bima":

		accountUser, _ := service.UserRepositoryInterface.GetUserAccountBimaByID(tx, userProfile.User.Id)
		if len(accountUser.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("user account paylater not found"), requestId, []string{"user account paylater not found"}, service.Logger, tx)
		}

		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Tabungan Bima"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
		orderEntity.PaymentCash = orderRequest.TotalBill

		desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, userProfile.User.IdDesa)
		if len(desa.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("desa account paylater not found"), requestId, []string{"desa account paylater not found"}, service.Logger, tx)
		}

		err = service.InveliAPIRepositoryInterface.ApiPayment(desa.NoRekening, accountUser.Code, userProfile.User.InveliAccessToken, orderRequest.TotalBill, 0)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}
	}

	// Create Order
	err = service.OrderRepositoryInterface.CreateOrder(tx, orderEntity)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order"}, service.Logger, tx)

	// Create order items
	err = service.OrderItemPpobRepositoryInterface.CreateOrderItemPpob(tx, orderItemsPpob)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order items"}, service.Logger, tx)

	if orderRequest.PaymentMethod == "tabungan_bima" || orderRequest.PaymentMethod == "paylater" {
		response := service.PrepaidPulsaTopup(requestId, orderRequest.CustomerId, orderEntity.RefId, orderItemsPpob.ProductCode)
		ppobDetailPrepaidPulsa.StatusTopUp = response.Data.Status
		ppobDetailPrepaidPulsa.TopupProccesingDate = null.NewTime(time.Now(), true)
		ppobDetailPrepaidPulsa.LastBalance = response.Data.Balance
	}

	err = service.PpobDetailRepositoryInterface.CreateOrderPpobDetailPrepaidPulsa(tx, ppobDetailPrepaidPulsa)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order items"}, service.Logger, tx)

	runtime.GOMAXPROCS(1)
	mssg := "Order Preapaid Pulsa Baru Dari " + userProfile.NamaLengkap + " ID Order " + orderEntity.NumberOrder + " VIA " + paymentChannel.Alias
	go service.SendMessageToTelegram(mssg, desa.ChatIdTelegram, desa.TokenBot)

	commit := tx.Commit()
	exceptions.PanicIfError(commit.Error, requestId, service.Logger)

	createOrderResponse = response.ToCreateOrderResponse(orderEntity, paymentChannel)
	return createOrderResponse
}

func (service *OrderServiceImplementation) CreateOrderPrepaidPln(requestId, idUser, idDesa, productType string, orderRequest *request.CreateOrderPrepaidRequest) (createOrderResponse response.CreateOrderResponse) {
	var err error

	request.ValidateRequest(service.Validate, orderRequest, requestId, service.Logger)

	// Get data user
	userProfile, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(userProfile.User.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user not found"), requestId, []string{"user not found"}, service.Logger)
	}

	tx := service.DB.Begin()
	// make object
	orderEntity := &entity.Order{}

	// Generate number and id order
	numberOrder := service.GenerateNumberOrder(idDesa)
	orderEntity.Id = utilities.RandomUUID()
	orderEntity.IdUser = idUser
	orderEntity.IdDesa = idDesa
	orderEntity.NumberOrder = numberOrder
	orderEntity.NamaLengkap = userProfile.NamaLengkap
	orderEntity.Email = userProfile.Email
	orderEntity.Phone = userProfile.User.Phone
	orderEntity.ProductType = productType
	orderEntity.PaymentPoint = orderRequest.PaymentPoint
	orderEntity.OrderedDate = time.Now()
	orderEntity.PaymentMethod = orderRequest.PaymentMethod
	orderEntity.PaymentChannel = orderRequest.PaymentChannel
	orderEntity.TotalBill = orderRequest.TotalBill
	orderEntity.PaymentFee = orderRequest.PaymentFee
	orderEntity.OrderType = 2
	orderEntity.RefId = utilities.GenerateRefId()

	// Create Request
	inquiryPlnData := service.OrderInquiryPrepaidPln(requestId, orderRequest.CustomerId)
	// Get Data from iak
	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + "pl"))
	body, _ := json.Marshal(map[string]interface{}{
		"status":   "all",
		"username": config.GetConfig().Ppob.Username,
		"sign":     hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	typePpob := "pln"
	urlString := config.GetConfig().Ppob.PrepaidHost + "/pricelist/" + typePpob + "/" + typePpob
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

	//  cek request

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	// fmt.Printf("body: %s\n", data)

	defer resp.Body.Close()

	priceLists := &ppob.PrepaidPriceListResponse{}
	if err = json.Unmarshal([]byte(data), &priceLists); err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	var totalHarga float64
	var product []string
	var qty []int
	var price []float64
	orderItemsPpob := &entity.OrderItemPpob{}
	ppobDetailPrepaidPln := &entity.PpobDetailPrepaidPln{}
	for _, priceList := range priceLists.Data.Data {
		if priceList.ProductCode == orderRequest.ProductCode {
			totalHarga = priceList.ProductPrice
			orderItemsPpob.Id = utilities.RandomUUID()
			orderItemsPpob.IdOrder = orderEntity.Id
			orderItemsPpob.IdUser = userProfile.IdUser
			orderItemsPpob.RefId = orderEntity.RefId
			orderItemsPpob.ProductCode = priceList.ProductCode
			orderItemsPpob.ProductType = productType
			orderItemsPpob.TotalTagihan = priceList.ProductPrice
			orderItemsPpob.IconUrl = priceList.IconUrl
			orderItemsPpob.CreatedAt = time.Now()
			orderItemsPpob.BillDetail = fmt.Sprintf("%s\n", data)

			ppobDetailPrepaidPln.Id = utilities.RandomUUID()
			ppobDetailPrepaidPln.IdOrderItemPpob = orderItemsPpob.Id
			ppobDetailPrepaidPln.ProductCode = priceList.ProductCode
			ppobDetailPrepaidPln.ProductName = priceList.ProductNominal
			ppobDetailPrepaidPln.ProductDescription = priceList.ProductDescription
			ppobDetailPrepaidPln.CustomerId = orderRequest.CustomerId
			ppobDetailPrepaidPln.MeterNo = inquiryPlnData.MeterNo
			ppobDetailPrepaidPln.SubscriberId = inquiryPlnData.SubscriberId
			ppobDetailPrepaidPln.CustomerName = inquiryPlnData.Name
			ppobDetailPrepaidPln.SegmentPower = inquiryPlnData.SegmentPower
			ppobDetailPrepaidPln.StatusTopUp = -1
			if orderRequest.PaymentMethod == "cc" {
				product = append(product, orderItemsPpob.ProductCode)
				qty = append(qty, 1)
				price = append(price, orderItemsPpob.TotalTagihan)
			}
			break
		}
	}

	fmt.Println("Total Harga = ", totalHarga)
	fmt.Println("Total Harga request = ", orderRequest.TotalBill)

	if ((totalHarga + 1500) + orderRequest.PaymentFee) != (orderRequest.TotalBill + orderRequest.PaymentFee + orderEntity.PaymentPoint) {
		exceptions.PanicIfErrorWithRollback(errors.New("harga tidak sama"), requestId, []string{"harga tidak sama"}, service.Logger, tx)
	}

	// Get detail payment channel
	paymentChannel, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(tx, orderRequest.PaymentChannel)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error get payment by code"}, service.Logger, tx)
	if len(paymentChannel.Id) == 0 {
		exceptions.PanicIfRecordNotFoundWithRollback(err, requestId, []string{"payment not found"}, service.Logger, tx)
	}

	// Get Desa
	desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, userProfile.User.IdDesa)
	if len(desa.Id) == 0 {
		exceptions.PanicIfErrorWithRollback(errors.New("desa account paylater not found"), requestId, []string{"desa account paylater not found"}, service.Logger, tx)
	}

	switch orderRequest.PaymentMethod {
	case "point":
		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Point"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
	case "trf":
		// buat nomor acak
		rand.Seed(time.Now().UnixNano())
		min := 111
		max := 299
		rand3Number := rand.Intn(max-min+1) + min

		min2 := 11
		max2 := 99
		rand2Number := rand.Intn(max2-min2+1) + min

		sisaPembagi := math.Mod(orderRequest.TotalBill, 1000)
		var Total float64

		if sisaPembagi < 100 {
			Total = orderRequest.TotalBill + float64(rand3Number)
		} else if sisaPembagi >= 100 {
			Total = orderRequest.TotalBill + float64(rand2Number)
		}

		orderEntity.OrderStatus = 0
		orderEntity.PaymentStatus = 0
		orderEntity.PaymentNo = paymentChannel.NoAccountBank
		orderEntity.PaymentName = paymentChannel.NamaPemilikBank
		orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
		orderEntity.PaymentCash = Total

	case "va", "qris":
		orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		res := service.PaymentServiceInterface.VaQrisPay(requestId,
			&payment.IpaymuQrisVaRequest{
				Name:           userProfile.NamaLengkap,
				Phone:          userProfile.User.Phone,
				Email:          userProfile.Email,
				Amount:         orderEntity.PaymentCash,
				ReferenceId:    numberOrder,
				PaymentMethod:  orderRequest.PaymentMethod,
				PaymentChannel: orderRequest.PaymentChannel,
			},
		)

		if res.Status != 200 {
			fmt.Println("LOG RESPONSE IPAYMU = ", res)
			exceptions.PanicIfErrorWithRollback(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger, tx)
		} else if res.Status == 200 {
			paymentDueDate, _ := time.Parse("2006-01-02 15:04:05", res.Data.Expired)
			orderEntity.PaymentStatus = 0
			orderEntity.TrxId = res.Data.TransactionId
			orderEntity.PaymentNo = res.Data.PaymentNo
			orderEntity.PaymentName = res.Data.PaymentName
			orderEntity.PaymentDueDate = null.NewTime(paymentDueDate, true)
			orderEntity.OrderStatus = 0
		}

	case "cc":
		// tambahkan ongkos kirim
		product = append(product, "Payment Fee")
		qty = append(qty, 1)
		price = append(price, orderRequest.PaymentFee)

		res := service.PaymentServiceInterface.CreditCardPay(requestId,
			&payment.IpaymuCreditCardRequest{
				Product:       product,
				Qty:           qty,
				Price:         price,
				ReferenceId:   numberOrder,
				BuyerName:     userProfile.NamaLengkap,
				BuyerEmail:    userProfile.Email,
				BuyerPhone:    userProfile.User.Phone,
				PaymentMethod: orderRequest.PaymentMethod,
			},
		)

		if res.Status != 200 {
			fmt.Println("LOG RESPONSE IPAYMU = ", res)
			exceptions.PanicIfErrorWithRollback(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger, tx)
		} else if res.Status == 200 {
			orderEntity.PaymentStatus = 0
			orderEntity.PaymentNo = res.Data.Url
			orderEntity.PaymentName = "Credit Card"
			orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
			orderEntity.OrderStatus = 0
			orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		}

	case "paylater":
		var isMerchant float64
		var totalAmount float64

		// Set Is Merchant 0
		isMerchant = 0

		// Validasi Saldo Bupda
		saldoBupda, err := service.InveliAPIRepositoryInterface.GetSaldoBupda(userProfile.User.InveliAccessToken, desa.GroupIdBupda)

		if err != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error saldo bupda "+err.Error()), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		if saldoBupda <= 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("saldo bupda kurang"), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		// Get Bunga
		bunga, errr := service.InveliAPIRepositoryInterface.GetLoanProduct(userProfile.User.InveliAccessToken)
		if errr != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error get loan product "+err.Error()), requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}

		// Get Loan Product
		loandProductID, errr := service.InveliAPIRepositoryInterface.GetLoanProductId(userProfile.User.InveliAccessToken)
		if errr != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error get loan product id "+err.Error()), requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}

		if len(loandProductID) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("loan product id not found"), requestId, []string{"loan product id not found"}, service.Logger, tx)
		}

		// Get Account User
		accountUser, _ := service.UserRepositoryInterface.GetUserAccountPaylaterByID(tx, userProfile.User.Id)
		if len(accountUser.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("user account paylater not found"), requestId, []string{"user account paylater not found"}, service.Logger, tx)
		}

		// Validasi Tunggakan Paylater
		// loanID, err := service.InveliAPIRepositoryInterface.GetRiwayatPinjaman(userProfile.User.InveliAccessToken, userProfile.User.InveliIDMember)
		// if err != nil {
		// 	log.Println("error get riwayat pinjaman", err.Error())
		// 	exceptions.PanicIfError(err, requestId, service.Logger)
		// }

		// if len(loanID) != 0 {
		// 	exceptions.PanicIfBadRequest(errors.New("masih ada tunggakan"), requestId, []string{"anda masih memiliki tunggakan"}, service.Logger)
		// }

		totalAmount = orderRequest.TotalBill + orderRequest.PaymentFee

		err = service.InveliAPIRepositoryInterface.InveliCreatePaylater(userProfile.User.InveliAccessToken, userProfile.User.InveliIDMember, accountUser.IdAccount, orderRequest.TotalBill, totalAmount, isMerchant, bunga, loandProductID, desa.NoRekening)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(errors.New("error care pinjaman "+err.Error()), requestId, []string{"Mohon maaf transaksi belum bisa dilakukan"}, service.Logger, tx)
		}

		// if time.Now().Local().Day() < 25 {
		// 	orderEntity.PaymentDueDate = null.NewTime(time.Date(time.Now().Year(), time.Now().Month(), 25, 0, 0, 0, 0, time.Local), true)
		// } else if time.Now().Local().Day() >= 25 {
		// 	orderEntity.PaymentDueDate = null.NewTime(time.Date(time.Now().Year(), time.Now().Month()+1, 25, 0, 0, 0, 0, time.Local), true)
		// }

		orderEntity.PaymentDueDate = null.NewTime(time.Now().AddDate(0, 0, 30), true)

		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Paylater"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
		orderEntity.PaymentCash = orderRequest.TotalBill + orderRequest.PaymentFee

		// err = service.InveliAPIRepositoryInterface.ApiPayment(desa.NoRekening, accountUser.Code, userProfile.User.InveliAccessToken, orderRequest.TotalBill, isMerchant)
		// if err != nil {
		// 	exceptions.PanicIfErrorWithRollback(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		// }

		var jmlOrder float64
		jmlOrderPayLate, err := service.OrderRepositoryInterface.FindOrderPayLaterById(service.DB, idUser)
		if err != nil {
			log.Println(err.Error())
		}
		jmlOrder = 0
		for _, v := range jmlOrderPayLate {
			jmlOrder = jmlOrder + v.TotalBill
		}

		userPaylaterFlag, _ := service.UserRepositoryInterface.GetUserPayLaterFlagThisMonth(service.DB, idUser)

		if (int(jmlOrder) + int(orderRequest.TotalBill)) > (userPaylaterFlag.TanggungRentengFlag * 1000000) {
			service.UserRepositoryInterface.UpdateUserPayLaterFlag(service.DB, idUser, &entity.UsersPaylaterFlag{
				TanggungRentengFlag: userPaylaterFlag.TanggungRentengFlag + 1,
			})
		}

	case "tabungan_bima":

		accountUser, _ := service.UserRepositoryInterface.GetUserAccountBimaByID(tx, userProfile.User.Id)
		if len(accountUser.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("user account paylater not found"), requestId, []string{"user account paylater not found"}, service.Logger, tx)
		}

		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Tabungan Bima"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
		orderEntity.PaymentCash = orderRequest.TotalBill

		desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, userProfile.User.IdDesa)
		if len(desa.Id) == 0 {
			exceptions.PanicIfErrorWithRollback(errors.New("desa account paylater not found"), requestId, []string{"desa account paylater not found"}, service.Logger, tx)
		}

		err = service.InveliAPIRepositoryInterface.ApiPayment(desa.NoRekening, accountUser.Code, userProfile.User.InveliAccessToken, orderRequest.TotalBill, 0)
		if err != nil {
			exceptions.PanicIfErrorWithRollback(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger, tx)
		}
	}

	// Create Order
	err = service.OrderRepositoryInterface.CreateOrder(tx, orderEntity)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order"}, service.Logger, tx)

	// Create order items
	err = service.OrderItemPpobRepositoryInterface.CreateOrderItemPpob(tx, orderItemsPpob)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order items"}, service.Logger, tx)

	if orderRequest.PaymentMethod == "tabungan_bima" || orderRequest.PaymentMethod == "paylater" {
		response := service.PrepaidPulsaTopup(requestId, orderRequest.CustomerId, orderEntity.RefId, orderItemsPpob.ProductCode)
		ppobDetailPrepaidPln.StatusTopUp = response.Data.Status
		ppobDetailPrepaidPln.TopupProccesingDate = null.NewTime(time.Now(), true)
		ppobDetailPrepaidPln.LastBalance = response.Data.Balance
	}

	err = service.PpobDetailRepositoryInterface.CreateOrderPpobDetailPrepaidPln(tx, ppobDetailPrepaidPln)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order items"}, service.Logger, tx)

	runtime.GOMAXPROCS(1)
	mssg := "Order Preapaid PLN Baru Dari " + userProfile.NamaLengkap + " ID Order " + orderEntity.NumberOrder + " VIA " + paymentChannel.Alias
	go service.SendMessageToTelegram(mssg, desa.ChatIdTelegram, desa.TokenBot)

	commit := tx.Commit()
	exceptions.PanicIfError(commit.Error, requestId, service.Logger)

	createOrderResponse = response.ToCreateOrderResponse(orderEntity, paymentChannel)
	return createOrderResponse
}

func (service *OrderServiceImplementation) CreateOrderSembako(requestId, idUser, idDesa string, accountType int, orderRequest *request.CreateOrderRequest) (createOrderResponse response.CreateOrderResponse) {
	var err error

	request.ValidateRequest(service.Validate, orderRequest, requestId, service.Logger)

	// Get data user
	userProfile, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(userProfile.User.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user not found"), requestId, []string{"user not found"}, service.Logger)
	}

	// Get data user shipping status
	userShippingAddress, _ := service.UserShippingAddressRepositoryInterface.FindUserShippingAddressByAddress(service.DB, orderRequest.AlamatPengiriman)

	// Get data user cart
	userCartItems, err := service.CartRepositoryInterface.FindCartByUser(service.DB, userProfile.User.Id)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(userCartItems) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("items in cart not found"), requestId, []string{"items in cart not found"}, service.Logger)
	}

	// make object
	orderEntity := &entity.Order{}

	// Generate number and id order
	numberOrder := service.GenerateNumberOrder(idDesa)
	orderEntity.Id = utilities.RandomUUID()
	orderEntity.IdUser = idUser
	orderEntity.IdDesa = idDesa
	orderEntity.NumberOrder = numberOrder
	orderEntity.NamaLengkap = userProfile.NamaLengkap
	orderEntity.Email = userProfile.Email
	orderEntity.Phone = userProfile.User.Phone
	orderEntity.AlamatPengiriman = orderRequest.AlamatPengiriman
	orderEntity.Catatan = orderRequest.CatatanKurir
	orderEntity.ShippingCost = orderRequest.ShippingCost
	orderEntity.ProductType = "sembako"
	orderEntity.PaymentPoint = orderRequest.PaymentPoint
	orderEntity.OrderedDate = time.Now()
	orderEntity.PaymentMethod = orderRequest.PaymentMethod
	orderEntity.PaymentChannel = orderRequest.PaymentChannel
	orderEntity.TotalBill = orderRequest.TotalBill + orderRequest.PaymentFee
	orderEntity.PaymentFee = orderRequest.PaymentFee
	orderEntity.Longitude = userShippingAddress.Longitude
	orderEntity.Latitude = userShippingAddress.Latitude

	orderEntity.OrderType = 1

	// create order items from cart
	var orderItems []entity.OrderItem
	var product []string
	var qty []int
	var price []float64
	var totalPrice float64
	orderItemsEntity := &entity.OrderItem{}
	for _, item := range userCartItems {
		orderItemsEntity.Id = utilities.RandomUUID()
		orderItemsEntity.IdOrder = orderEntity.Id
		orderItemsEntity.IdProductDesa = item.IdProductDesa
		orderItemsEntity.IdUser = userProfile.User.Id
		orderItemsEntity.ProductName = item.ProductsDesa.ProductsMaster.ProductName
		orderItemsEntity.PictureUrl = item.ProductsDesa.ProductsMaster.PictureUrl
		orderItemsEntity.Thumbnail = item.ProductsDesa.ProductsMaster.Thumbnail
		orderItemsEntity.FlagPromo = item.ProductsDesa.IsPromo
		orderItemsEntity.Description = item.ProductsDesa.ProductsMaster.Description
		orderItemsEntity.Qty = item.Qty
		if accountType == 1 {
			orderItemsEntity.Price = item.ProductsDesa.Price
			if item.ProductsDesa.IsPromo == 1 {
				orderItemsEntity.PriceAfterDiscount = item.ProductsDesa.PricePromo
				orderItemsEntity.TotalPrice = float64(orderItemsEntity.Qty) * item.ProductsDesa.PricePromo
			} else if item.ProductsDesa.IsPromo == 0 {
				orderItemsEntity.PriceAfterDiscount = 0
				orderItemsEntity.TotalPrice = float64(orderItemsEntity.Qty) * orderItemsEntity.Price
			}
		} else if accountType == 2 {
			orderItemsEntity.Price = item.ProductsDesa.ProductsMaster.PriceGrosir
			orderItemsEntity.TotalPrice = float64(orderItemsEntity.Qty) * orderItemsEntity.Price
		}
		orderItemsEntity.CreatedAt = time.Now()

		totalPrice = totalPrice + orderItemsEntity.TotalPrice
		orderItems = append(orderItems, *orderItemsEntity)
		if orderRequest.PaymentMethod == "cc" {
			product = append(product, orderItemsEntity.ProductName)
			qty = append(qty, orderItemsEntity.Qty)
			price = append(price, orderItemsEntity.Price)
		}
	}
	orderEntity.SubTotal = totalPrice

	// Checking total bill from FE
	log.Println("Harga kalkulasi server 1 = ", totalPrice+orderRequest.ShippingCost)
	log.Println("Harga dari client 1 = ", orderRequest.TotalBill+orderRequest.PaymentPoint)
	if (totalPrice + orderRequest.ShippingCost) != (orderRequest.TotalBill + orderRequest.PaymentPoint) {
		exceptions.PanicIfBadRequest(errors.New("harga tidak sama"), requestId, []string{"harga tidak sama"}, service.Logger)
	}

	log.Println("Harga kalkulasi server 2 = ", totalPrice+orderRequest.ShippingCost+orderRequest.PaymentFee)
	log.Println("Harga dari client 2 = ", (orderRequest.TotalBill+orderRequest.PaymentFee)+orderRequest.PaymentPoint)
	// Checking total payment from FE
	if (totalPrice + orderRequest.ShippingCost + orderRequest.PaymentFee) != ((orderRequest.TotalBill + orderRequest.PaymentFee) + orderRequest.PaymentPoint) {
		exceptions.PanicIfBadRequest(errors.New("harga tidak sama dengan payment cash"), requestId, []string{"harga tidak sama dengan payment cash"}, service.Logger)
	}

	// Get detail payment channel
	paymentChannel, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(service.DB, orderRequest.PaymentChannel)
	if err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	if len(paymentChannel.Id) == 0 {
		exceptions.PanicIfRecordNotFound(err, requestId, []string{"payment not found"}, service.Logger)
	}

	// Get Desa
	desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, userProfile.User.IdDesa)
	if len(desa.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("desa account paylater not found"), requestId, []string{"desa account paylater not found"}, service.Logger)
	}

	switch orderRequest.PaymentMethod {
	case "cod":
		orderEntity.OrderStatus = 1
		orderEntity.PaymentCash = orderRequest.TotalBill
		orderEntity.PaymentName = "Cash On Delivery"
	case "point":
		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Point"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
	case "trf":
		// buat nomor acak
		rand.Seed(time.Now().UnixNano())
		min := 111
		max := 299
		rand3Number := rand.Intn(max-min+1) + min

		min2 := 11
		max2 := 99
		rand2Number := rand.Intn(max2-min2+1) + min

		sisaPembagi := math.Mod(orderRequest.TotalBill, 1000)
		var Total float64

		if sisaPembagi < 100 {
			Total = orderRequest.TotalBill + float64(rand3Number)
		} else if sisaPembagi >= 100 {
			Total = orderRequest.TotalBill + float64(rand2Number)
		}

		orderEntity.OrderStatus = 0
		orderEntity.PaymentStatus = 0
		orderEntity.PaymentNo = paymentChannel.NoAccountBank
		orderEntity.PaymentName = paymentChannel.NamaPemilikBank
		orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
		orderEntity.PaymentCash = Total

	case "va", "qris":
		orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		res := service.PaymentServiceInterface.VaQrisPay(requestId,
			&payment.IpaymuQrisVaRequest{
				Name:           userProfile.NamaLengkap,
				Phone:          userProfile.User.Phone,
				Email:          userProfile.Email,
				Amount:         orderEntity.PaymentCash,
				ReferenceId:    numberOrder,
				PaymentMethod:  orderRequest.PaymentMethod,
				PaymentChannel: orderRequest.PaymentChannel,
			},
		)

		if res.Status != 200 {
			fmt.Println("LOG RESPONSE IPAYMU = ", res)
			exceptions.PanicIfBadRequest(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger)
		} else if res.Status == 200 {
			paymentDueDate, _ := time.Parse("2006-01-02 15:04:05", res.Data.Expired)
			orderEntity.PaymentStatus = 0
			orderEntity.TrxId = res.Data.TransactionId
			orderEntity.PaymentNo = res.Data.PaymentNo
			orderEntity.PaymentName = res.Data.PaymentName
			orderEntity.PaymentDueDate = null.NewTime(paymentDueDate, true)
			orderEntity.OrderStatus = 0
		}

	case "cc":
		// tambahkan ongkos kirim
		product = append(product, "Shipping Cost", "Payment Fee")
		qty = append(qty, 1, 1)
		price = append(price, orderRequest.ShippingCost, orderRequest.PaymentFee)

		res := service.PaymentServiceInterface.CreditCardPay(requestId,
			&payment.IpaymuCreditCardRequest{
				Product:       product,
				Qty:           qty,
				Price:         price,
				ReferenceId:   numberOrder,
				BuyerName:     userProfile.NamaLengkap,
				BuyerEmail:    userProfile.Email,
				BuyerPhone:    userProfile.User.Phone,
				PaymentMethod: orderRequest.PaymentMethod,
			},
		)

		if res.Status != 200 {
			fmt.Println("LOG RESPONSE IPAYMU = ", res)
			exceptions.PanicIfBadRequest(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger)
		} else if res.Status == 200 {
			orderEntity.PaymentStatus = 0
			orderEntity.PaymentNo = res.Data.Url
			orderEntity.PaymentName = "Credit Card"
			orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
			orderEntity.OrderStatus = 0
			orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		}

	case "paylater":

		orderPaylater := service.PaymentServiceInterface.PayWithPaylater(userProfile.User.InveliAccessToken, userProfile.User.InveliIDMember, desa.GroupIdBupda, desa.NoRekening, userProfile.User.Id, orderRequest.TotalBill, orderRequest.PaymentFee)

		orderEntity.PaymentDueDate = orderPaylater.PaymentDueDate
		orderEntity.OrderStatus = orderPaylater.OrderStatus
		orderEntity.PaymentStatus = orderPaylater.PaymentStatus
		orderEntity.PaymentName = orderPaylater.PaymentName
		orderEntity.PaymentSuccessDate = orderPaylater.PaymentSuccessDate
		orderEntity.PaymentCash = orderPaylater.PaymentCash

	case "tabungan_bima":

		accountUser, _ := service.UserRepositoryInterface.GetUserAccountBimaByID(service.DB, userProfile.User.Id)
		if len(accountUser.Id) == 0 {
			exceptions.PanicIfRecordNotFound(errors.New("user account paylater not found"), requestId, []string{"user account paylater not found"}, service.Logger)
		}

		orderEntity.OrderStatus = 1
		orderEntity.PaymentStatus = 1
		orderEntity.PaymentName = "Tabungan Bima"
		orderEntity.PaymentSuccessDate = null.NewTime(time.Now(), true)
		orderEntity.PaymentCash = orderRequest.TotalBill

		desa, _ := service.DesaRepositoryInterface.FindDesaById(service.DB, userProfile.User.IdDesa)
		if len(desa.Id) == 0 {
			exceptions.PanicIfRecordNotFound(errors.New("desa account paylater not found"), requestId, []string{"desa account paylater not found"}, service.Logger)
		}

		err = service.InveliAPIRepositoryInterface.ApiPayment(desa.NoRekening, accountUser.Code, userProfile.User.InveliAccessToken, orderRequest.TotalBill, 0)
		if err != nil {
			exceptions.PanicIfBadRequest(err, requestId, []string{strings.TrimPrefix(err.Error(), "graphql: ")}, service.Logger)
		}
	}

	// // Get Desa
	runtime.GOMAXPROCS(1)
	mssg := "Order Sembako Baru Dari " + userProfile.NamaLengkap + " ID Order " + orderEntity.NumberOrder + " VIA " + paymentChannel.Alias
	go service.SendMessageToTelegram(mssg, desa.ChatIdTelegram, desa.TokenBot)

	// Create Order
	tx := service.DB.Begin()
	err = service.OrderRepositoryInterface.CreateOrder(tx, orderEntity)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order"}, service.Logger, tx)

	// Create order items
	err = service.OrderItemRepositoryInterface.CreateOrderItem(tx, orderItems)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order items"}, service.Logger, tx)

	// Delete items in cart
	err = service.CartRepositoryInterface.DeleteCartByUser(tx, idUser, userCartItems)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error delete items in cart"}, service.Logger, tx)

	// update stock jika payment methodnya point
	if orderRequest.PaymentMethod == "point" || orderRequest.PaymentMethod == "cod" || orderRequest.PaymentMethod == "tabungan_bima" || orderRequest.PaymentMethod == "paylater" {
		service.ProductDesaServiceInterface.UpdateProductStock(requestId, orderEntity.Id, tx)
	}

	commit := tx.Commit()
	exceptions.PanicIfError(commit.Error, requestId, service.Logger)

	createOrderResponse = response.ToCreateOrderResponse(orderEntity, paymentChannel)
	return createOrderResponse
}

func (service *OrderServiceImplementation) SendMessageToTelegram(message, chatId, token string) {
	url, _ := url.Parse("https://api.telegram.org/bot" + token + "/sendMessage?chat_id=" + chatId + "&text=" + message)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("an error occured %v", err)
	}
	defer resp.Body.Close()
}

func (service *OrderServiceImplementation) FindOrderByUser(requestId, idUser string, orderStatus int) (orderResponses []response.FindOrderByUserResponse) {
	orders, err := service.OrderRepositoryInterface.FindOrderByUser(service.DB, idUser, orderStatus)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orders) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}
	orderResponses = response.ToFindOrderByUserResponse(orders)
	return orderResponses
}

func (service *OrderServiceImplementation) FindOrderSembakoById(requestId, idOrder string) (orderResponse response.FindOrderSembakoByIdResponse) {
	var err error

	// Get order by id order
	order, err := service.OrderRepositoryInterface.FindOrderById(service.DB, idOrder)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(order.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}

	// Get order items by id oder
	orderItems, err := service.OrderItemRepositoryInterface.FindOrderItemsByIdOrder(service.DB, idOrder)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orderItems) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order item not found"), requestId, []string{"order item not found"}, service.Logger)
	}

	// Payment
	payment, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(service.DB, order.PaymentChannel)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(payment.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("payment not found"), requestId, []string{"order item not found"}, service.Logger)
	}

	orderResponse = response.ToFindOrderSembakoByIdResponse(order, orderItems, payment)
	return orderResponse
}

func (service *OrderServiceImplementation) FindOrderPaymentById(requestId, idOrder string) (orderResponse response.OrderPayment) {
	var err error

	// Get order by id order
	order, err := service.OrderRepositoryInterface.FindOrderById(service.DB, idOrder)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(order.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}

	// Payment
	payment, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(service.DB, order.PaymentChannel)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(payment.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("payment not found"), requestId, []string{"order item not found"}, service.Logger)
	}

	orderResponse = response.ToFindOrderPaymentyIdResponse(order, payment)
	return orderResponse
}

func (service *OrderServiceImplementation) FindOrderPrepaidPulsaById(requestId, idOrder string, productType string) (orderResponse response.FindOrderPrepaidPulsaByIdResponse) {
	var err error

	// Get order by id order
	order, err := service.OrderRepositoryInterface.FindOrderPrepaidPulsaById(service.DB, idOrder, productType)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(order.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}

	// order items ppob
	orderItemsPpob, err := service.OrderItemPpobRepositoryInterface.FindOrderItemsPpobByIdOrder(service.DB, order.Id)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orderItemsPpob.Id) == 0 {
		exceptions.PanicIfBadRequest(errors.New("order items ppob not found"), requestId, []string{"order items ppob not found"}, service.Logger)
	}

	// Get detail prepaid pulsa
	ppobDetailPrepaidPulsa, err := service.PpobDetailRepositoryInterface.FindPpobDetailPrepaidPulsaById(service.DB, orderItemsPpob.Id)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orderItemsPpob.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order item ppb not found"), requestId, []string{"order item ppob not found"}, service.Logger)
	}

	// Payment
	payment, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(service.DB, order.PaymentChannel)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(payment.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("payment not found"), requestId, []string{"order item not found"}, service.Logger)
	}

	orderResponse = response.ToFindOrderPrepaidPulsaByIdResponse(order, orderItemsPpob, ppobDetailPrepaidPulsa, payment)
	return orderResponse
}

func (service *OrderServiceImplementation) FindOrderPrepaidPlnById(requestId, idOrder string) (orderResponse response.FindOrderPrepaidPlnByIdResponse) {
	var err error

	// Get order by id order
	order, err := service.OrderRepositoryInterface.FindOrderById(service.DB, idOrder)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(order.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}

	// order items ppob
	orderItemsPpob, err := service.OrderItemPpobRepositoryInterface.FindOrderItemsPpobByIdOrder(service.DB, order.Id)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orderItemsPpob.Id) == 0 {
		exceptions.PanicIfBadRequest(errors.New("order items ppob not found"), requestId, []string{"order items ppob not found"}, service.Logger)
	}

	// Get detail prepaid pulsa
	ppobDetailPrepaidPln, err := service.PpobDetailRepositoryInterface.FindPpobDetailPrepaidPlnById(service.DB, orderItemsPpob.Id)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orderItemsPpob.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order item ppb not found"), requestId, []string{"order item ppob not found"}, service.Logger)
	}

	// Payment
	payment, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(service.DB, order.PaymentChannel)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(payment.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("payment not found"), requestId, []string{"order item not found"}, service.Logger)
	}

	orderResponse = response.ToFindOrderPrepaidPlnByIdResponse(order, orderItemsPpob, ppobDetailPrepaidPln, payment)
	return orderResponse
}

func (service *OrderServiceImplementation) FindOrderPostpaidPlnById(requestId, idOrder string) (orderResponse response.FindOrderPostpaidPlnByIdResponse) {
	var err error

	// Get order by id order
	order, err := service.OrderRepositoryInterface.FindOrderById(service.DB, idOrder)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(order.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}

	// order items ppob
	orderItemsPpob, err := service.OrderItemPpobRepositoryInterface.FindOrderItemsPpobByIdOrder(service.DB, order.Id)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orderItemsPpob.Id) == 0 {
		exceptions.PanicIfBadRequest(errors.New("order items ppob not found"), requestId, []string{"order items ppob not found"}, service.Logger)
	}

	// Get detail prepaid pulsa
	ppobDetailPostpaidPln, err := service.PpobDetailRepositoryInterface.FindPpobDetailPostpaidPlnById(service.DB, orderItemsPpob.Id)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orderItemsPpob.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order item ppb not found"), requestId, []string{"order item ppob not found"}, service.Logger)
	}

	postpaidPlnDetail := []ppob.InquiryPostpaidPlnDetail{}

	if err = json.Unmarshal([]byte(ppobDetailPostpaidPln.JsonDetailTagihan), &postpaidPlnDetail); err != nil {
		exceptions.PanicIfBadRequest(errors.New("INVALID DATA"), requestId, []string{"INVALID DATA"}, service.Logger)
	}

	// Payment
	payment, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(service.DB, order.PaymentChannel)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(payment.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("payment not found"), requestId, []string{"order item not found"}, service.Logger)
	}

	orderResponse = response.ToFindOrderPostpaidPlnByIdResponse(order, orderItemsPpob, ppobDetailPostpaidPln, payment, postpaidPlnDetail)
	return orderResponse
}

func (service *OrderServiceImplementation) FindOrderPostpaidPdamById(requestId, idOrder string) (orderResponse response.FindOrderPostpaidPdamByIdResponse) {
	var err error

	// Get order by id order
	order, err := service.OrderRepositoryInterface.FindOrderById(service.DB, idOrder)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(order.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}

	// order items ppob
	orderItemsPpob, err := service.OrderItemPpobRepositoryInterface.FindOrderItemsPpobByIdOrder(service.DB, order.Id)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orderItemsPpob.Id) == 0 {
		exceptions.PanicIfBadRequest(errors.New("order items ppob not found"), requestId, []string{"order items ppob not found"}, service.Logger)
	}

	// Get detail prepaid pulsa
	ppobDetailPostpaidPdam, err := service.PpobDetailRepositoryInterface.FindPpobDetailPostpaidPdamById(service.DB, orderItemsPpob.Id)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orderItemsPpob.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order item ppb not found"), requestId, []string{"order item ppob not found"}, service.Logger)
	}

	postpaidPdamDetail := []ppob.InquiryPostpaidPdamBillDetail{}

	if err = json.Unmarshal([]byte(ppobDetailPostpaidPdam.JsonDetailTagihan), &postpaidPdamDetail); err != nil {
		exceptions.PanicIfBadRequest(errors.New("INVALID DATA"), requestId, []string{"INVALID DATA"}, service.Logger)
	}

	// Payment
	payment, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(service.DB, order.PaymentChannel)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(payment.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("payment not found"), requestId, []string{"order item not found"}, service.Logger)
	}

	orderResponse = response.ToFindOrderPostpaidPdamByIdResponse(order, orderItemsPpob, ppobDetailPostpaidPdam, payment, postpaidPdamDetail)
	return orderResponse
}

func (service *OrderServiceImplementation) CancelOrderById(requestId string, orderRequest *request.OrderIdRequest) {
	var err error
	//Get order detail
	request.ValidateRequest(service.Validate, orderRequest, requestId, service.Logger)
	order, err := service.OrderRepositoryInterface.FindOrderById(service.DB, orderRequest.IdOrder)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(order.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}

	// Update status order
	err = service.OrderRepositoryInterface.UpdateOrderByIdOrder(service.DB, orderRequest.IdOrder, &entity.Order{
		OrderStatus:       9,
		PaymentStatus:     9,
		OrderCanceledDate: null.NewTime(time.Now(), true),
	})
	exceptions.PanicIfError(err, requestId, service.Logger)
}

func (service *OrderServiceImplementation) CompleteOrderById(requestId string, orderRequest *request.OrderIdRequest) {
	var err error
	//Get order detail
	request.ValidateRequest(service.Validate, orderRequest, requestId, service.Logger)
	order, err := service.OrderRepositoryInterface.FindOrderById(service.DB, orderRequest.IdOrder)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(order.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}

	tx := service.DB.Begin()

	// Check apakah order sudah sampai ditujuan
	if order.OrderStatus != 4 {
		exceptions.PanicIfBadRequest(errors.New("orderan belum sampai di tujuan"), requestId, []string{"orderan blum sampai di tujuan"}, service.Logger)
	}

	// Update status order menjadi selesai
	err = service.OrderRepositoryInterface.UpdateOrderByIdOrder(tx, orderRequest.IdOrder, &entity.Order{
		OrderStatus:        5,
		OrderCompletedDate: null.NewTime(time.Now(), true),
	})
	exceptions.PanicIfError(err, requestId, service.Logger)

	// Script untuk bonus point

	commit := tx.Commit()
	exceptions.PanicIfError(commit.Error, requestId, service.Logger)
}

func (service *OrderServiceImplementation) UpdatePaymentStatusOrder(requestId string, updatePaymentStatusOrderRequest *request.UpdatePaymentStatusOrderRequest) {
	var err error
	//Get order detail
	request.ValidateRequest(service.Validate, updatePaymentStatusOrderRequest, requestId, service.Logger)

	order, err := service.OrderRepositoryInterface.FindOrderByNumberOrder(service.DB, updatePaymentStatusOrderRequest.ReferenceId)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(order.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}

	// cek apakah order status masih 0
	if order.OrderStatus != 0 {
		exceptions.PanicIfBadRequest(errors.New("order tidak dalam status 0"), requestId, []string{"order tidak dalam status 0"}, service.Logger)
	}

	// Check status order ke ipaymu

	paymentStatus := service.PaymentServiceInterface.CheckPaymentStatus(requestId, updatePaymentStatusOrderRequest.TrxId)

	if paymentStatus.Status == 1 || paymentStatus.Status == 6 {
		if order.OrderType == 1 {
			err = service.OrderRepositoryInterface.UpdateOrderByIdOrder(service.DB, order.Id, &entity.Order{
				OrderStatus:        1,
				PaymentStatus:      1,
				PaymentSuccessDate: null.NewTime(time.Now(), true),
			})
			exceptions.PanicIfError(err, requestId, service.Logger)

			service.ProductDesaServiceInterface.UpdateProductStock(requestId, order.Id, service.DB)

			// Update stock
		} else if order.OrderType == 2 {
			err = service.OrderRepositoryInterface.UpdateOrderByIdOrder(service.DB, order.Id, &entity.Order{
				OrderStatus:        2,
				PaymentStatus:      1,
				PaymentSuccessDate: null.NewTime(time.Now(), true),
			})
			exceptions.PanicIfError(err, requestId, service.Logger)

			orderItemsPpob, err := service.OrderItemPpobRepositoryInterface.FindOrderItemsPpobByIdOrder(service.DB, order.Id)
			exceptions.PanicIfError(err, requestId, service.Logger)
			if len(orderItemsPpob.Id) == 0 {
				exceptions.PanicIfBadRequest(errors.New("order items ppob not found"), requestId, []string{"order items ppob not found"}, service.Logger)
			}

			switch order.ProductType {
			case "prepaid_pulsa", "prepaid_data":
				ppobDetailPrepaidPulsa, err := service.PpobDetailRepositoryInterface.FindPpobDetailPrepaidPulsaById(service.DB, orderItemsPpob.Id)
				exceptions.PanicIfError(err, requestId, service.Logger)
				if len(ppobDetailPrepaidPulsa.Id) == 0 {
					exceptions.PanicIfBadRequest(errors.New("ppob detail prepaid pulsa not found"), requestId, []string{"ppob detail prepaid pulsa not found"}, service.Logger)
				}

				response := service.PrepaidPulsaTopup(requestId, ppobDetailPrepaidPulsa.CustomerId, order.RefId, orderItemsPpob.ProductCode)

				err = service.PpobDetailRepositoryInterface.UpdatePpobPrepaidPulsaById(service.DB, ppobDetailPrepaidPulsa.Id, &entity.PpobDetailPrepaidPulsa{
					StatusTopUp:         response.Data.Status,
					TopupProccesingDate: null.NewTime(time.Now(), true),
					LastBalance:         response.Data.Balance,
				})
				exceptions.PanicIfError(err, requestId, service.Logger)
			case "prepaid_pln":
				ppobDetailPrepaidPulsa, err := service.PpobDetailRepositoryInterface.FindPpobDetailPrepaidPlnById(service.DB, orderItemsPpob.Id)
				exceptions.PanicIfError(err, requestId, service.Logger)
				if len(ppobDetailPrepaidPulsa.Id) == 0 {
					exceptions.PanicIfBadRequest(errors.New("ppob detail prepaid pulsa not found"), requestId, []string{"ppob detail prepaid pulsa not found"}, service.Logger)
				}

				response := service.PrepaidPulsaTopup(requestId, ppobDetailPrepaidPulsa.CustomerId, order.RefId, orderItemsPpob.ProductCode)

				err = service.PpobDetailRepositoryInterface.UpdatePpobPrepaidPlnById(service.DB, ppobDetailPrepaidPulsa.Id, &entity.PpobDetailPrepaidPln{
					StatusTopUp:         response.Data.Status,
					NoToken:             response.Data.Sn,
					TopupProccesingDate: null.NewTime(time.Now(), true),
					LastBalance:         response.Data.Balance,
				})
				exceptions.PanicIfError(err, requestId, service.Logger)
			case "postpaid_pln":
				ppobDetailPostpaidPln, err := service.PpobDetailRepositoryInterface.FindPpobDetailPostpaidPlnById(service.DB, orderItemsPpob.Id)
				exceptions.PanicIfError(err, requestId, service.Logger)
				if len(ppobDetailPostpaidPln.Id) == 0 {
					exceptions.PanicIfBadRequest(errors.New("ppob detail postpaid pln not found"), requestId, []string{"ppob detail postpaid pln not found"}, service.Logger)
				}

				response := service.PostpaidTopupPln(requestId, ppobDetailPostpaidPln.CustomerId, ppobDetailPostpaidPln.OrderItemPpob.TrId, orderItemsPpob.ProductCode)

				err = service.PpobDetailRepositoryInterface.UpdatePpobPostpaidPlnById(service.DB, ppobDetailPostpaidPln.Id, &entity.PpobDetailPostpaidPln{
					StatusTopUp:         3,
					TopupProccesingDate: null.NewTime(time.Now(), true),
					LastBalance:         response.Balance,
				})
				exceptions.PanicIfError(err, requestId, service.Logger)

			case "postpaid_pdam":
				ppobDetailPostpaidPdam, err := service.PpobDetailRepositoryInterface.FindPpobDetailPostpaidPdamById(service.DB, orderItemsPpob.Id)
				exceptions.PanicIfError(err, requestId, service.Logger)
				if len(ppobDetailPostpaidPdam.Id) == 0 {
					exceptions.PanicIfBadRequest(errors.New("ppob detail postpaid pdam not found"), requestId, []string{"ppob detail postpaid pdam not found"}, service.Logger)
				}

				response := service.PostpaidTopupPdam(requestId, ppobDetailPostpaidPdam.CustomerId, ppobDetailPostpaidPdam.OrderItemPpob.TrId, orderItemsPpob.ProductCode)

				err = service.PpobDetailRepositoryInterface.UpdatePpobPostpaidPdamById(service.DB, ppobDetailPostpaidPdam.Id, &entity.PpobDetailPostpaidPdam{
					StatusTopUp:         3,
					TopupProccesingDate: null.NewTime(time.Now(), true),
					LastBalance:         response.Balance,
				})
				exceptions.PanicIfError(err, requestId, service.Logger)
			}

		} else {
			exceptions.PanicIfBadRequest(errors.New("order type not found"), requestId, []string{"order type not found"}, service.Logger)
		}
	} else {
		exceptions.PanicIfBadRequest(errors.New("status pembayaran belum terbayar"), requestId, []string{"status pembayaran belum terbayar"}, service.Logger)
	}
}

func (service *OrderServiceImplementation) PrepaidPulsaTopup(requestId string, customerId, refId, productCode string) *ppob.TopupPrepaidPulsaResponse {
	var err error

	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + refId))
	body, _ := json.Marshal(map[string]interface{}{
		"username":     config.GetConfig().Ppob.Username,
		"ref_id":       refId,
		"customer_id":  customerId,
		"product_code": productCode,
		"sign":         hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PrepaidHost + "/top-up"

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

	defer resp.Body.Close()

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	topupPrepaidPulsaResponse := &ppob.TopupPrepaidPulsaResponse{}

	if err = json.Unmarshal([]byte(data), &topupPrepaidPulsaResponse); err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	return topupPrepaidPulsaResponse
}

func (service *OrderServiceImplementation) PostpaidTopupPln(requestId string, customerId string, TrxId int, productCode string) *ppob.TopupPostaidPlnDataResponse {
	var err error

	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + string(rune(TrxId))))
	body, _ := json.Marshal(map[string]interface{}{
		"commands": "pay-pasca",
		"username": config.GetConfig().Ppob.Username,
		"tr_id":    TrxId,
		"sign":     hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PostpaidUrl

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

	defer resp.Body.Close()

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	topupPostpaidPlnResonse := &ppob.TopupPostaidPlnResponse{}

	if err = json.Unmarshal([]byte(data), &topupPostpaidPlnResonse); err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	return &topupPostpaidPlnResonse.Data
}

func (service *OrderServiceImplementation) PostpaidTopupPdam(requestId string, customerId string, TrxId int, productCode string) *ppob.TopupPostaidPdamDataResponse {
	var err error

	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + string(rune(TrxId))))
	body, _ := json.Marshal(map[string]interface{}{
		"commands": "pay-pasca",
		"username": config.GetConfig().Ppob.Username,
		"tr_id":    TrxId,
		"sign":     hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PostpaidUrl

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

	defer resp.Body.Close()

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	topupPostpaidPdamResonse := &ppob.TopupPostaidPdamResponse{}

	if err = json.Unmarshal([]byte(data), &topupPostpaidPdamResonse); err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	return &topupPostpaidPdamResonse.Data
}

func (service *OrderServiceImplementation) OrderInquiryPrepaidPln(requestId string, customerId string) (inquiryPrepaidPlnResponse response.InquiryPrepaidPlnResponse) {
	var err error

	// Create Request
	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + customerId))
	body, _ := json.Marshal(map[string]interface{}{
		"username":    config.GetConfig().Ppob.Username,
		"customer_id": customerId,
		"sign":        hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PrepaidHost + "/inquiry-pln"

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

	defer resp.Body.Close()

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	inquiryPrepaidPln := &ppob.InquiryPrepaidPln{}

	if err = json.Unmarshal([]byte(data), &inquiryPrepaidPln); err != nil {
		inquiryPrepaidPlnErrorHandle := &ppob.InquiryPrepaidPlnErrorHandle{}
		if err = json.Unmarshal([]byte(data), &inquiryPrepaidPlnErrorHandle); err != nil {
			exceptions.PanicIfError(err, requestId, service.Logger)
		} else {
			if inquiryPrepaidPlnErrorHandle.Data.Rc == "208" {
				exceptions.PanicIfBadRequest(errors.New("INVALID DATA"), requestId, []string{"INVALID DATA"}, service.Logger)
			}
			if inquiryPrepaidPlnErrorHandle.Data.Rc == "14" {
				exceptions.PanicIfBadRequest(errors.New("costumer id not found"), requestId, []string{"Costumer Id Not Found"}, service.Logger)
			}
			exceptions.PanicIfError(err, requestId, service.Logger)
		}
	}

	if inquiryPrepaidPln.Data.Rc != "00" {
		fmt.Printf("body: %s\n", inquiryPrepaidPln.Data)
		exceptions.PanicIfError(errors.New("error from IAK"), requestId, service.Logger)
	}

	inquiryPrepaidPlnResponse = response.ToInquiryPrepaidPlnResponse(inquiryPrepaidPln)

	return inquiryPrepaidPlnResponse
}

func (service *OrderServiceImplementation) CallbackPpobTransaction(requestId string, ppobCallbackRequest *request.PpobCallbackRequest) {
	var err error

	fmt.Println("request = ", ppobCallbackRequest)

	order, err := service.OrderRepositoryInterface.FindOrderByRefId(service.DB, ppobCallbackRequest.Data.RefId)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(order.RefId) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}

	signCheck := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + string(order.RefId)))
	log.Println("sign check = ", hex.EncodeToString(signCheck[:]))
	log.Println("sign = ", ppobCallbackRequest.Data.Sign)

	// cek sign dari iak dengan signcheck
	if hex.EncodeToString(signCheck[:]) != ppobCallbackRequest.Data.Sign {
		exceptions.PanicIfBadRequest(errors.New("sign not match"), requestId, []string{"sign not match"}, service.Logger)
	}

	orderItemsPpob, err := service.OrderItemPpobRepositoryInterface.FindOrderItemsPpobByIdOrder(service.DB, order.Id)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(orderItemsPpob.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order item ppob not found"), requestId, []string{"order item ppob not found"}, service.Logger)
	}

	switch order.ProductType {
	case "prepaid_pulsa":
		detailPrepaidPulsa, err := service.PpobDetailRepositoryInterface.FindPpobDetailPrepaidPulsaById(service.DB, orderItemsPpob.Id)
		exceptions.PanicIfError(err, requestId, service.Logger)
		if len(orderItemsPpob.Id) == 0 {
			exceptions.PanicIfRecordNotFound(errors.New("detail prepaid pulsa not found"), requestId, []string{"detail prepaid pulsa not found"}, service.Logger)
		}

		// transaksi sukses
		if ppobCallbackRequest.Data.Status == "1" {
			// update order
			err = service.OrderRepositoryInterface.UpdateOrderByIdOrder(service.DB, order.Id, &entity.Order{
				OrderStatus:        5,
				OrderCompletedDate: null.NewTime(time.Now(), true),
			})
			exceptions.PanicIfError(err, requestId, service.Logger)

			balance, _ := strconv.ParseFloat(ppobCallbackRequest.Data.Balance, 32)

			err = service.PpobDetailRepositoryInterface.UpdatePpobPrepaidPulsaById(service.DB, detailPrepaidPulsa.Id, &entity.PpobDetailPrepaidPulsa{
				StatusTopUp:      1,
				TopupSuccessDate: null.NewTime(time.Now(), true),
				LastBalance:      balance,
			})
			exceptions.PanicIfError(err, requestId, service.Logger)
		} else if ppobCallbackRequest.Data.Status == "2" {
			// Transaksi failed
			err = service.OrderRepositoryInterface.UpdateOrderByIdOrder(service.DB, order.Id, &entity.Order{
				OrderStatus:       9,
				OrderCanceledDate: null.NewTime(time.Now(), true),
			})
			exceptions.PanicIfError(err, requestId, service.Logger)

			balance, _ := strconv.ParseFloat(ppobCallbackRequest.Data.Balance, 32)

			err = service.PpobDetailRepositoryInterface.UpdatePpobPrepaidPulsaById(service.DB, detailPrepaidPulsa.Id, &entity.PpobDetailPrepaidPulsa{
				StatusTopUp:     2,
				TopupFailedDate: null.NewTime(time.Now(), true),
				LastBalance:     balance,
			})
			exceptions.PanicIfError(err, requestId, service.Logger)

		} else {
			exceptions.PanicIfBadRequest(errors.New("status not found"), requestId, []string{"error string"}, service.Logger)
		}

	case "prepaid_pln":
		detailPrepaidPln, err := service.PpobDetailRepositoryInterface.FindPpobDetailPrepaidPlnById(service.DB, orderItemsPpob.Id)
		exceptions.PanicIfError(err, requestId, service.Logger)
		if len(orderItemsPpob.Id) == 0 {
			exceptions.PanicIfRecordNotFound(errors.New("detail prepaid pln not found"), requestId, []string{"detail prepaid pln not found"}, service.Logger)
		}

		if ppobCallbackRequest.Data.Status == "1" {
			// update order
			err = service.OrderRepositoryInterface.UpdateOrderByIdOrder(service.DB, order.Id, &entity.Order{
				OrderStatus:        5,
				OrderCompletedDate: null.NewTime(time.Now(), true),
			})
			exceptions.PanicIfError(err, requestId, service.Logger)

			balance, _ := strconv.ParseFloat(ppobCallbackRequest.Data.Balance, 32)

			err = service.PpobDetailRepositoryInterface.UpdatePpobPrepaidPlnById(service.DB, detailPrepaidPln.Id, &entity.PpobDetailPrepaidPln{
				StatusTopUp:      1,
				TopupSuccessDate: null.NewTime(time.Now(), true),
				LastBalance:      balance,
				NoToken:          ppobCallbackRequest.Data.Sn,
			})
			exceptions.PanicIfError(err, requestId, service.Logger)
		} else if ppobCallbackRequest.Data.Status == "2" {
			err = service.OrderRepositoryInterface.UpdateOrderByIdOrder(service.DB, order.Id, &entity.Order{
				OrderStatus:       9,
				OrderCanceledDate: null.NewTime(time.Now(), true),
			})
			exceptions.PanicIfError(err, requestId, service.Logger)

			balance, _ := strconv.ParseFloat(ppobCallbackRequest.Data.Balance, 32)

			err = service.PpobDetailRepositoryInterface.UpdatePpobPrepaidPlnById(service.DB, detailPrepaidPln.Id, &entity.PpobDetailPrepaidPln{
				StatusTopUp:     2,
				TopupFailedDate: null.NewTime(time.Now(), true),
				LastBalance:     balance,
			})
			exceptions.PanicIfError(err, requestId, service.Logger)
		} else {
			exceptions.PanicIfBadRequest(errors.New("status not found"), requestId, []string{"error string"}, service.Logger)
		}

	case "postpaid_pln":
		detailPostpaidPln, err := service.PpobDetailRepositoryInterface.FindPpobDetailPostpaidPlnById(service.DB, orderItemsPpob.Id)
		exceptions.PanicIfError(err, requestId, service.Logger)
		if len(orderItemsPpob.Id) == 0 {
			exceptions.PanicIfRecordNotFound(errors.New("detail postpaid pln not found"), requestId, []string{"detail postpaid pln not found"}, service.Logger)
		}

		if ppobCallbackRequest.Data.Status == "1" {
			// update order
			err = service.OrderRepositoryInterface.UpdateOrderByIdOrder(service.DB, order.Id, &entity.Order{
				OrderStatus:        5,
				OrderCompletedDate: null.NewTime(time.Now(), true),
			})
			exceptions.PanicIfError(err, requestId, service.Logger)

			balance, _ := strconv.ParseFloat(ppobCallbackRequest.Data.Balance, 32)

			err = service.PpobDetailRepositoryInterface.UpdatePpobPostpaidPlnById(service.DB, detailPostpaidPln.Id, &entity.PpobDetailPostpaidPln{
				StatusTopUp:      1,
				TopupSuccessDate: null.NewTime(time.Now(), true),
				LastBalance:      balance,
			})
			exceptions.PanicIfError(err, requestId, service.Logger)
		} else if ppobCallbackRequest.Data.Status == "2" {
			err = service.OrderRepositoryInterface.UpdateOrderByIdOrder(service.DB, order.Id, &entity.Order{
				OrderStatus:       9,
				OrderCanceledDate: null.NewTime(time.Now(), true),
			})
			exceptions.PanicIfError(err, requestId, service.Logger)

			balance, _ := strconv.ParseFloat(ppobCallbackRequest.Data.Balance, 32)

			err = service.PpobDetailRepositoryInterface.UpdatePpobPostpaidPlnById(service.DB, detailPostpaidPln.Id, &entity.PpobDetailPostpaidPln{
				StatusTopUp:     2,
				TopupFailedDate: null.NewTime(time.Now(), true),
				LastBalance:     balance,
			})
			exceptions.PanicIfError(err, requestId, service.Logger)
		} else {
			exceptions.PanicIfBadRequest(errors.New("status not found"), requestId, []string{"error string"}, service.Logger)
		}

	case "postpaid_pdam":
		detailPostpaidPdam, err := service.PpobDetailRepositoryInterface.FindPpobDetailPostpaidPdamById(service.DB, orderItemsPpob.Id)
		exceptions.PanicIfError(err, requestId, service.Logger)
		if len(orderItemsPpob.Id) == 0 {
			exceptions.PanicIfRecordNotFound(errors.New("detail postpaid pln not found"), requestId, []string{"detail postpaid pln not found"}, service.Logger)
		}

		if ppobCallbackRequest.Data.Status == "1" {
			// update order
			err = service.OrderRepositoryInterface.UpdateOrderByIdOrder(service.DB, order.Id, &entity.Order{
				OrderStatus:        5,
				OrderCompletedDate: null.NewTime(time.Now(), true),
			})
			exceptions.PanicIfError(err, requestId, service.Logger)

			balance, _ := strconv.ParseFloat(ppobCallbackRequest.Data.Balance, 32)

			err = service.PpobDetailRepositoryInterface.UpdatePpobPostpaidPdamById(service.DB, detailPostpaidPdam.Id, &entity.PpobDetailPostpaidPdam{
				StatusTopUp:      1,
				TopupSuccessDate: null.NewTime(time.Now(), true),
				LastBalance:      balance,
			})
			exceptions.PanicIfError(err, requestId, service.Logger)
		} else if ppobCallbackRequest.Data.Status == "2" {
			err = service.OrderRepositoryInterface.UpdateOrderByIdOrder(service.DB, order.Id, &entity.Order{
				OrderStatus:       9,
				OrderCanceledDate: null.NewTime(time.Now(), true),
			})
			exceptions.PanicIfError(err, requestId, service.Logger)

			balance, _ := strconv.ParseFloat(ppobCallbackRequest.Data.Balance, 32)

			err = service.PpobDetailRepositoryInterface.UpdatePpobPostpaidPdamById(service.DB, detailPostpaidPdam.Id, &entity.PpobDetailPostpaidPdam{
				StatusTopUp:     2,
				TopupFailedDate: null.NewTime(time.Now(), true),
				LastBalance:     balance,
			})
			exceptions.PanicIfError(err, requestId, service.Logger)
		} else {
			exceptions.PanicIfBadRequest(errors.New("status not found"), requestId, []string{"error string"}, service.Logger)
		}
	}
}

package service

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/payment"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"github.com/tensuqiuwulu/be-service-bupda-bali/utilities"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type OrderServiceInterface interface {
	CreateOrder(requestId, idUser, idDesa string, accountType int, orderRequest *request.CreateOrderRequest) (createOrderResponse response.CreateOrderResponse)
	FindOrderByUser(requestId, idUser string, orderStatus int) (orderResponses []response.FindOrderByUserResponse)
	FindOrderById(requestId, idOrder string) (orderResponse response.FindOrderByIdResponse)
	CancelOrderById(requestId string, orderRequest *request.OrderIdRequest)
	CompleteOrderById(requestId string, orderRequest *request.OrderIdRequest)
	UpdatePaymentStatusOrderById(requestId string, orderRequest *request.OrderIdRequest)
}

type OrderServiceImplementation struct {
	DB                                *gorm.DB
	Validate                          *validator.Validate
	Logger                            *logrus.Logger
	OrderRepositoryInterface          repository.OrderRepositoryInterface
	UserRepositoryInterface           repository.UserRepositoryInterface
	PaymentServiceInterface           PaymentServiceInterface
	CartRepositoryInterface           repository.CartRepositoryInterface
	OrderItemRepositoryInterface      repository.OrderItemRepositoryInterface
	PaymentChannelRepositoryInterface repository.PaymentChannelRepositoryInterface
	ProductDesaRepositoryInterface    repository.ProductDesaRepositoryInterface
	ProductDesaServiceInterface       ProductDesaServiceInterface
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
) OrderServiceInterface {
	return &OrderServiceImplementation{
		DB:                                db,
		Validate:                          validate,
		Logger:                            logger,
		OrderRepositoryInterface:          orderRepositoryInterface,
		UserRepositoryInterface:           userRepositoryInterface,
		PaymentServiceInterface:           paymentServiceInterface,
		CartRepositoryInterface:           cartRepositoryInterface,
		OrderItemRepositoryInterface:      orderItemRepositoryInterface,
		PaymentChannelRepositoryInterface: paymentChannelRepositoryInterface,
		ProductDesaRepositoryInterface:    productDesaRepositoryInterface,
		ProductDesaServiceInterface:       productDesaServiceInterface,
	}
}

func (service *OrderServiceImplementation) GenerateNumberOrder() (numberOrder string) {
	now := time.Now()
	for {
		rand.Seed(time.Now().UTC().UnixNano())
		generateCode := 100000 + rand.Intn(999999-100000)
		numberOrder = "ORDER/" + now.Format("20060102") + "/" + fmt.Sprint(generateCode)

		// Check number order if exist
		order, _ := service.OrderRepositoryInterface.FindOrderByNumberOrder(service.DB, numberOrder)
		if len(order.Id) == 0 {
			break
		}
	}
	return numberOrder
}

func (service *OrderServiceImplementation) CreateOrder(requestId, idUser, idDesa string, accountType int, orderRequest *request.CreateOrderRequest) (createOrderResponse response.CreateOrderResponse) {
	var err error

	request.ValidateRequest(service.Validate, orderRequest, requestId, service.Logger)

	// Get data user
	userProfile, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(userProfile.User.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user not found"), requestId, []string{"user not found"}, service.Logger)
	}

	// Get data user cart
	userCartItems, err := service.CartRepositoryInterface.FindCartByUser(service.DB, userProfile.User.Id)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(userCartItems) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("items in cart not found"), requestId, []string{"items in cart not found"}, service.Logger)
	}

	tx := service.DB.Begin()
	// make object
	orderEntity := &entity.Order{}

	// Generate number and id order
	numberOrder := service.GenerateNumberOrder()
	orderEntity.Id = utilities.RandomUUID()
	orderEntity.IdUser = idUser
	orderEntity.NumberOrder = numberOrder
	orderEntity.NamaLengkap = userProfile.NamaLengkap
	orderEntity.Email = userProfile.Email
	orderEntity.Phone = userProfile.User.Phone
	orderEntity.AlamatPengiriman = orderRequest.AlamatPengiriman
	orderEntity.Catatan = orderRequest.CatatanKurir
	orderEntity.ShippingCost = orderRequest.ShippingCost
	orderEntity.OrderName = "Order Product"
	// if orderRequest.PaymentMethod != "trf" && {
	// 	orderEntity.PaymentCash = orderRequest.PaymentCash + orderEntity.PaymentFee
	// }
	orderEntity.PaymentPoint = orderRequest.PaymentPoint
	orderEntity.OrderedDate = time.Now()
	orderEntity.PaymentMethod = orderRequest.PaymentMethod
	orderEntity.PaymentChannel = orderRequest.PaymentChannel
	orderEntity.TotalBill = orderRequest.TotalBill
	orderEntity.PaymentFee = orderRequest.PaymentFee

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
			orderItemsEntity.Price = item.ProductsDesa.ProductsMaster.Price
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
	if (totalPrice + orderRequest.ShippingCost) != (orderRequest.TotalBill + orderRequest.PaymentPoint) {
		exceptions.PanicIfErrorWithRollback(errors.New("harga tidak sama"), requestId, []string{"harga tidak sama"}, service.Logger, tx)
	}

	// Checking total payment from FE
	if (totalPrice + orderRequest.ShippingCost + orderRequest.PaymentFee) != ((orderRequest.TotalBill + orderRequest.PaymentFee) + orderRequest.PaymentPoint) {
		exceptions.PanicIfErrorWithRollback(errors.New("harga tidak sama dengan payment cash"), requestId, []string{"harga tidak sama dengan payment cash"}, service.Logger, tx)
	}

	// Get detail payment channel
	paymentChannel, err := service.PaymentChannelRepositoryInterface.FindPaymentChannelByCode(tx, orderRequest.PaymentChannel)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error get payment by code"}, service.Logger, tx)
	if len(paymentChannel.Id) == 0 {
		exceptions.PanicIfRecordNotFoundWithRollback(err, requestId, []string{"payment not found"}, service.Logger, tx)
	}

	switch orderRequest.PaymentMethod {
	case "cod":
		orderEntity.OrderStatus = 0
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
			exceptions.PanicIfErrorWithRollback(errors.New("error response ipaymu"), requestId, []string{"Error response ipaymu"}, service.Logger, tx)
		} else if res.Status == 200 {
			orderEntity.PaymentStatus = 0
			orderEntity.PaymentNo = res.Data.Url
			orderEntity.PaymentName = "Credit Card"
			orderEntity.PaymentDueDate = null.NewTime(time.Now().Add(time.Hour*24), true)
			orderEntity.OrderStatus = 0
			orderEntity.PaymentCash = orderRequest.TotalBill + orderEntity.PaymentFee
		}
	}

	// Create Order
	err = service.OrderRepositoryInterface.CreateOrder(tx, orderEntity)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order"}, service.Logger, tx)

	// Create order items
	err = service.OrderItemRepositoryInterface.CreateOrderItem(tx, orderItems)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error create order items"}, service.Logger, tx)

	// Delete items in cart
	err = service.CartRepositoryInterface.DeleteCartByUser(tx, idUser, userCartItems)
	exceptions.PanicIfErrorWithRollback(err, requestId, []string{"error delete items in cart"}, service.Logger, tx)

	// update stock jika payment methodnya point
	if orderRequest.PaymentMethod == "point" {
		service.ProductDesaServiceInterface.UpdateProductStock(requestId, orderEntity.Id, tx)
	}

	commit := tx.Commit()
	exceptions.PanicIfError(commit.Error, requestId, service.Logger)

	createOrderResponse = response.ToCreateOrderResponse(orderEntity, paymentChannel)
	return createOrderResponse
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

func (service *OrderServiceImplementation) FindOrderById(requestId, idOrder string) (orderResponse response.FindOrderByIdResponse) {
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

	orderResponse = response.ToFindOrderByIdResponse(order, orderItems)
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
}

func (service *OrderServiceImplementation) UpdatePaymentStatusOrderById(requestId string, orderRequest *request.OrderIdRequest) {
	var err error
	//Get order detail
	request.ValidateRequest(service.Validate, orderRequest, requestId, service.Logger)
	order, err := service.OrderRepositoryInterface.FindOrderById(service.DB, orderRequest.IdOrder)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(order.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("order not found"), requestId, []string{"order not found"}, service.Logger)
	}

	// cek apakah order status masih 0
	if order.OrderStatus != 0 {
		exceptions.PanicIfBadRequest(errors.New("order tidak dalam status 0"), requestId, []string{"order tidak dalam status 0"}, service.Logger)
	}

	// Check status order ke ipaymu

}

package service

import (
	"errors"
	"time"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"github.com/tensuqiuwulu/be-service-bupda-bali/utilities"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type CartServiceInterface interface {
	CreateCart(requestId string, idUser string, createCartRequest *request.CreateCartRequest) string
	UpdateCart(requestId string, updateCartRequest *request.UpdateCartRequest) string
	FindCartByUser(requestId string, idUser string, accountType int, idDesa string) (cartResponses response.FindCartByIdUserResponse)
}

type CartServiceImplementation struct {
	DB                          *gorm.DB
	Validate                    *validator.Validate
	Logger                      *logrus.Logger
	CartRepositoryInterface     repository.CartRepositoryInterface
	ProductDesaServiceInterface repository.ProductDesaRepositoryInterface
	SettingRepositoryInterface  repository.SettingRepositoryInterface
	DesaRepositoryInterface     repository.DesaRepositoryInterface
}

func NewCartService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	cartRepositoryInterface repository.CartRepositoryInterface,
	productDesaRepositoryInterface repository.ProductDesaRepositoryInterface,
	settingRepositoryInterface repository.SettingRepositoryInterface,
	desaRepositoryInterface repository.DesaRepositoryInterface,
) CartServiceInterface {
	return &CartServiceImplementation{
		DB:                          db,
		Validate:                    validate,
		Logger:                      logger,
		CartRepositoryInterface:     cartRepositoryInterface,
		ProductDesaServiceInterface: productDesaRepositoryInterface,
		SettingRepositoryInterface:  settingRepositoryInterface,
		DesaRepositoryInterface:     desaRepositoryInterface,
	}
}

func (service *CartServiceImplementation) CreateCart(requestId string, idUser string, createcCartRequest *request.CreateCartRequest) string {
	var err error

	request.ValidateRequest(service.Validate, createcCartRequest, requestId, service.Logger)

	// Check product if exist in cart
	cartResult, err := service.CartRepositoryInterface.FindCartByProductDesa(service.DB, idUser, createcCartRequest.IdProduct)
	exceptions.PanicIfError(err, requestId, service.Logger)

	// Check stock
	productResult, err := service.ProductDesaServiceInterface.FindProductDesaById(service.DB, createcCartRequest.IdProduct)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if createcCartRequest.Qty > productResult.StockOpname {
		exceptions.PanicIfBadRequest(errors.New("out of stock"), requestId, []string{"stock opname tidak cukup"}, service.Logger)
	}

	if len(cartResult.Id) == 0 {
		// add product to cart
		cartEntity := &entity.Cart{
			Id:            utilities.RandomUUID(),
			IdUser:        idUser,
			IdProductDesa: createcCartRequest.IdProduct,
			Qty:           createcCartRequest.Qty,
			CreatedDate:   time.Now(),
		}

		err := service.CartRepositoryInterface.CreateCart(service.DB, cartEntity)
		exceptions.PanicIfError(err, requestId, service.Logger)

		return cartEntity.Id
	} else if len(cartResult.Id) != 0 {
		// update product in cart
		cartEntity := &entity.Cart{
			Qty:         cartResult.Qty + createcCartRequest.Qty,
			UpdatedDate: null.NewTime(time.Now(), true),
		}

		err := service.CartRepositoryInterface.UpdateCart(service.DB, cartResult.Id, cartEntity)
		exceptions.PanicIfError(err, requestId, service.Logger)
		return cartResult.Id
	} else {
		exceptions.PanicIfBadRequest(errors.New("error add product to cart"), requestId, []string{"error add product to cart"}, service.Logger)
		return ""
	}
}

func (service *CartServiceImplementation) UpdateCart(requestId string, updateCartRequest *request.UpdateCartRequest) string {
	var err error

	request.ValidateRequest(service.Validate, updateCartRequest, requestId, service.Logger)

	// Check product if exist in cart
	cartResult, err := service.CartRepositoryInterface.FindCartById(service.DB, updateCartRequest.IdCart)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(cartResult.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("cart not found"), requestId, []string{"cart not found"}, service.Logger)
	}

	if updateCartRequest.Qty <= 0 {
		err := service.CartRepositoryInterface.DeleteCartById(service.DB, updateCartRequest.IdCart)
		exceptions.PanicIfError(err, requestId, service.Logger)
	} else {
		cartEntity := &entity.Cart{
			Qty:         updateCartRequest.Qty,
			UpdatedDate: null.NewTime(time.Now(), true),
		}

		err := service.CartRepositoryInterface.UpdateCart(service.DB, cartResult.Id, cartEntity)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}
	return cartResult.Id
}

func (service *CartServiceImplementation) FindCartByUser(requestid string, idUser string, accountType int, idDesa string) (cartResponses response.FindCartByIdUserResponse) {
	carts, err := service.CartRepositoryInterface.FindCartByUser(service.DB, idUser)
	exceptions.PanicIfError(err, requestid, service.Logger)
	if len(carts) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("not found"), requestid, []string{"data not found"}, service.Logger)
	}

	desa, err := service.DesaRepositoryInterface.FindDesaById(service.DB, idDesa)
	exceptions.PanicIfError(err, requestid, service.Logger)
	if desa.Ongkir == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("shipping cost not found"), requestid, []string{"shipping cost not found"}, service.Logger)
	}

	cartResponses = response.ToFindCartByUserResponse(carts, desa.Ongkir, accountType)
	return cartResponses
}

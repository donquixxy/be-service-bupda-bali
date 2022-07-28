package service

import (
	"errors"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"github.com/tensuqiuwulu/be-service-bupda-bali/utilities"
	"gorm.io/gorm"
)

type UserShippingAddressServiceInterface interface {
	FindUserShippingAddressByIdUser(requestId string, idUser string) (userAShippingddressResponses []response.FindUserShippingAddress)
	CreateUserShippingAddress(requestId string, idUser string, userShippingAddressRequest *request.CreateUserShippingAddressRequest)
	DeleteUserShippingAddress(requestId string, userShippingAddressRequest *request.DeleteUserShippingAddressRequest)
}

type UserShippingAddressServiceImplementation struct {
	DB                                     *gorm.DB
	Validate                               *validator.Validate
	Logger                                 *logrus.Logger
	UserShippingAddressRepositoryInterface repository.UserShippingAddressRepositoryInterface
}

func NewUserShippingAddressService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	userShippingAddressRepositoryInterface repository.UserShippingAddressRepositoryInterface,
) UserShippingAddressServiceInterface {
	return &UserShippingAddressServiceImplementation{
		DB:                                     db,
		Validate:                               validate,
		Logger:                                 logger,
		UserShippingAddressRepositoryInterface: userShippingAddressRepositoryInterface,
	}
}

func (service *UserShippingAddressServiceImplementation) DeleteUserShippingAddress(requestId string, userShippingAddressRequest *request.DeleteUserShippingAddressRequest) {
	var err error
	// validate
	request.ValidateRequest(service.Validate, userShippingAddressRequest, requestId, service.Logger)

	shippingAddress, err := service.UserShippingAddressRepositoryInterface.FindUserShippingAddressById(service.DB, userShippingAddressRequest.IdUserShippingAddress)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(shippingAddress.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("not found"), requestId, []string{"addres not found"}, service.Logger)
	}

	err = service.UserShippingAddressRepositoryInterface.DeleteUserShippingAddress(service.DB, userShippingAddressRequest.IdUserShippingAddress)
	exceptions.PanicIfError(err, requestId, service.Logger)
}

func (service *UserShippingAddressServiceImplementation) CreateUserShippingAddress(requestId string, idUser string, createUserShippingAddressRequest *request.CreateUserShippingAddressRequest) {
	// validate request
	request.ValidateRequest(service.Validate, createUserShippingAddressRequest, requestId, service.Logger)

	userShippingAddressEntity := &entity.UserShippingAddress{}
	userShippingAddressEntity.Id = utilities.RandomUUID()
	userShippingAddressEntity.IdUser = idUser
	userShippingAddressEntity.AlamatPengiriman = createUserShippingAddressRequest.AlamatPengiriman
	userShippingAddressEntity.Latitude = createUserShippingAddressRequest.Latitude
	userShippingAddressEntity.Longitude = createUserShippingAddressRequest.Longitude
	userShippingAddressEntity.Radius = createUserShippingAddressRequest.Radius
	userShippingAddressEntity.StatusPrimary = createUserShippingAddressRequest.StatusPrimary
	userShippingAddressEntity.Catatan = createUserShippingAddressRequest.Catatan
	_, err := service.UserShippingAddressRepositoryInterface.CreateUserShippingAddress(service.DB, userShippingAddressEntity)
	exceptions.PanicIfError(err, requestId, service.Logger)
}

func (service *UserShippingAddressServiceImplementation) FindUserShippingAddressByIdUser(requestId string, idUser string) (userShippingAddressResponses []response.FindUserShippingAddress) {
	userShippingAddresss, err := service.UserShippingAddressRepositoryInterface.FindUserShippingAddressByIdUser(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(userShippingAddresss) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user address not found"), requestId, []string{"user address not found"}, service.Logger)
	}
	userShippingAddressResponses = response.ToFindUserShippingAddressResponse(userShippingAddresss)
	return userShippingAddressResponses
}

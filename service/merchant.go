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
	"gorm.io/gorm"
)

type MerchantServiceInterface interface {
	CreateMerchantApproveList(requestId, idUser, idDesa string, createcMerchantApproveListRequest *request.CreateMerchantApproveListRequest) string
	FindMerchantStatusApproveByUser(requestId, idUser string) (merchantResponse *response.FindMerchantStatusApproveByUserResponse)
}

type MerchantServiceImplementation struct {
	DB                             *gorm.DB
	Validate                       *validator.Validate
	Logger                         *logrus.Logger
	UserProfileRepositoryInterface repository.UserProfileRepositoryInterface
	MerchantRepositoryInterface    repository.MerchantRepositoryInterface
}

func NewMerchantService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	userProfileRepositoryInterface repository.UserProfileRepositoryInterface,
	merchantRepositoryInterface repository.MerchantRepositoryInterface,
) MerchantServiceInterface {
	return &MerchantServiceImplementation{
		DB:                             db,
		Validate:                       validate,
		Logger:                         logger,
		UserProfileRepositoryInterface: userProfileRepositoryInterface,
		MerchantRepositoryInterface:    merchantRepositoryInterface,
	}
}

func (service *MerchantServiceImplementation) FindMerchantStatusApproveByUser(requestId, idUser string) (merchantResponse *response.FindMerchantStatusApproveByUserResponse) {
	merchantApproveStatus, err := service.MerchantRepositoryInterface.FindMerchantStatusApproveByUser(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(merchantApproveStatus.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("status approve merchant not found"), requestId, []string{"status approve merchant not found"}, service.Logger)
	}

	response := response.ToFindMerchantStatusApproveByUserResponse(merchantApproveStatus)
	return &response
}

func (service *MerchantServiceImplementation) CreateMerchantApproveList(requestId, idUser, idDesa string, createcMerchantApproveListRequest *request.CreateMerchantApproveListRequest) string {
	var err error

	request.ValidateRequest(service.Validate, createcMerchantApproveListRequest, requestId, service.Logger)

	// Get User Profile
	userProfile, err := service.UserProfileRepositoryInterface.FindUserProfileByIdUser(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(userProfile.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("user profile not found"), requestId, []string{"user profile not found"}, service.Logger)
	}

	// add product to cart
	merchantApproveListEntity := &entity.MerchantApproveList{
		Id:            utilities.RandomUUID(),
		IdUser:        idUser,
		IdDesa:        idDesa,
		NamaLengkap:   userProfile.NamaLengkap,
		ApproveStatus: 1,
		MerchantName:  createcMerchantApproveListRequest.NamaMerchant,
		CreatedAt:     time.Now(),
	}

	err = service.MerchantRepositoryInterface.CreateMerchantApproveList(service.DB, merchantApproveListEntity)
	exceptions.PanicIfError(err, requestId, service.Logger)

	return ""
}

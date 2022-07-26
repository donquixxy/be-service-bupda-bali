package service

import (
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"gorm.io/gorm"
)

type PromoServiceInterface interface {
	FindPromo(requestId string, IdDesa string) (promosResponses []response.FindPromoResponse)
}

type PromoServiceImplementation struct {
	DB                       *gorm.DB
	Validate                 *validator.Validate
	Logger                   *logrus.Logger
	PromoRepositoryInterface repository.PromoRepositoryInterface
}

func NewPromoService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	promoServiceInterface repository.PromoRepositoryInterface,
) PromoServiceInterface {
	return &PromoServiceImplementation{
		DB:                       db,
		Validate:                 validate,
		Logger:                   logger,
		PromoRepositoryInterface: promoServiceInterface,
	}
}

func (service *PromoServiceImplementation) FindPromo(requestId string, IdDesa string) (promosResponses []response.FindPromoResponse) {
	promos, err := service.PromoRepositoryInterface.FindPromo(service.DB, IdDesa)
	exceptions.PanicIfRecordNotFound(err, requestId, []string{"Data not found"}, service.Logger)
	promosResponses = response.ToFindPromoResponse(promos)
	return promosResponses
}

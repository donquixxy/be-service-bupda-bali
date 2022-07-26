package service

import (
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"gorm.io/gorm"
)

type SettingServiceInterface interface {
	FindSettingShippingCost(requestId string, idDesa string) (settingResponse response.FindSettingShippingCostResponse)
}

type SettingServiceImplementation struct {
	DB                         *gorm.DB
	Validate                   *validator.Validate
	Logger                     *logrus.Logger
	SettingRepositoryInterface repository.SettingRepositoryInterface
}

func NewSettingService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	settingServiceInterface repository.SettingRepositoryInterface,
) SettingServiceInterface {
	return &SettingServiceImplementation{
		DB:                         db,
		Validate:                   validate,
		Logger:                     logger,
		SettingRepositoryInterface: settingServiceInterface,
	}
}

func (service *SettingServiceImplementation) FindSettingShippingCost(requestId string, idDesa string) (settingResponse response.FindSettingShippingCostResponse) {
	shippingCost, err := service.SettingRepositoryInterface.FindSettingShippingCost(service.DB, idDesa)
	exceptions.PanicIfRecordNotFound(err, requestId, []string{"Data not found"}, service.Logger)
	settingResponse = response.ToFindSettingShippingCostResponse(shippingCost)
	return settingResponse
}

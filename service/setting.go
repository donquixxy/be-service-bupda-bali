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
	FindNewVersion(requestId string, os int) (settingResponse response.FindVersionResponse)
}

type SettingServiceImplementation struct {
	DB                            *gorm.DB
	Validate                      *validator.Validate
	Logger                        *logrus.Logger
	SettingRepositoryInterface    repository.SettingRepositoryInterface
	DesaRepositoryInterface       repository.DesaRepositoryInterface
	AppVersionRepositoryInterface repository.AppVersionRepositoryInterface
}

func NewSettingService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	settingServiceInterface repository.SettingRepositoryInterface,
	desaRepositoryInterface repository.DesaRepositoryInterface,
	appVersionRepositoryInterface repository.AppVersionRepositoryInterface,
) SettingServiceInterface {
	return &SettingServiceImplementation{
		DB:                            db,
		Validate:                      validate,
		Logger:                        logger,
		SettingRepositoryInterface:    settingServiceInterface,
		DesaRepositoryInterface:       desaRepositoryInterface,
		AppVersionRepositoryInterface: appVersionRepositoryInterface,
	}
}

func (service *SettingServiceImplementation) FindNewVersion(requestId string, os int) (settingResponse response.FindVersionResponse) {
	var osType string
	if os == 1 {
		osType = "android"
	} else if os == 2 {
		osType = "ios"
	}

	version, err := service.AppVersionRepositoryInterface.FindVersionByOS(service.DB, osType)
	exceptions.PanicIfRecordNotFound(err, requestId, []string{"Data not found"}, service.Logger)
	settingResponse = response.ToFindNewVersionResponse(version, os)
	return settingResponse
}

func (service *SettingServiceImplementation) FindSettingShippingCost(requestId string, idDesa string) (settingResponse response.FindSettingShippingCostResponse) {
	shippingCost, _ := service.SettingRepositoryInterface.FindSettingShippingCost(service.DB, idDesa)
	desa, err := service.DesaRepositoryInterface.FindDesaById(service.DB, idDesa)
	if err != nil {
		exceptions.PanicIfRecordNotFound(err, requestId, []string{"error get shipping cost : ", err.Error()}, service.Logger)
	}

	if len(desa.Id) == 0 {
		exceptions.PanicIfRecordNotFound(err, requestId, []string{"shipping cost not found : ", err.Error()}, service.Logger)
	}

	settingResponse = response.ToFindSettingShippingCostResponse(shippingCost, desa.Ongkir)
	return settingResponse
}

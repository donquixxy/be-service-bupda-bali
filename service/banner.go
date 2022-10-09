package service

import (
	"errors"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"gorm.io/gorm"
)

type BannerServiceInterface interface {
	FindBannerByDesa(requestId, idDesa string) (bannersResponses []response.FindBannerByIdDesaResponse)
	FindBannerAll(requestId string) (bannersResponses []response.FindBannerAllResponse)
}

type BannerServiceImplementation struct {
	DB                           *gorm.DB
	Validate                     *validator.Validate
	Logger                       *logrus.Logger
	BannerReprepositoryInterface repository.BannerRepositoryInterface
}

func NewBannerService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	bannerReprepositoryInterface repository.BannerRepositoryInterface,
) BannerServiceInterface {
	return &BannerServiceImplementation{
		DB:                           db,
		Validate:                     validate,
		Logger:                       logger,
		BannerReprepositoryInterface: bannerReprepositoryInterface,
	}
}

func (service *BannerServiceImplementation) FindBannerByDesa(requestId, idDesa string) (bannersResponses []response.FindBannerByIdDesaResponse) {
	banners, err := service.BannerReprepositoryInterface.FindBannerByDesa(service.DB, idDesa)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(banners) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("not found"), requestId, []string{"data not found"}, service.Logger)
	}
	bannersResponses = response.ToFindBannerByIdDesaResponse(banners)
	return bannersResponses
}

func (service *BannerServiceImplementation) FindBannerAll(requestId string) (bannersResponses []response.FindBannerAllResponse) {
	banners, err := service.BannerReprepositoryInterface.FindBannerAll(service.DB)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(banners) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("not found"), requestId, []string{"data not found"}, service.Logger)
	}
	bannersResponses = response.ToFindBannerAllResponse(banners)
	return bannersResponses
}

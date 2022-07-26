package service

import (
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"gorm.io/gorm"
)

type KelurahanServiceInterface interface {
	FindKelurahanByIdKeca(requestId string, idKeca int) (kelurahanResponses []response.FindKelurahanByIdKecaResponse)
}

type KelurahanServiceImplementation struct {
	DB                           *gorm.DB
	Validate                     *validator.Validate
	Logger                       *logrus.Logger
	KelurahanRepositoryInterface repository.KelurahanRepositoryInterface
}

func NewKelurahanService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	kelurahanServiceInterface repository.KelurahanRepositoryInterface,
) KelurahanServiceInterface {
	return &KelurahanServiceImplementation{
		DB:                           db,
		Validate:                     validate,
		Logger:                       logger,
		KelurahanRepositoryInterface: kelurahanServiceInterface,
	}
}

func (service *KelurahanServiceImplementation) FindKelurahanByIdKeca(requestId string, idKeca int) (kelurahansResponses []response.FindKelurahanByIdKecaResponse) {
	kelurahans, err := service.KelurahanRepositoryInterface.FindKelurahanByIdKeca(service.DB, idKeca)
	exceptions.PanicIfRecordNotFound(err, requestId, []string{"Data not found"}, service.Logger)
	kelurahansResponses = response.ToFindKelurahanByIdKecaResponse(kelurahans)
	return kelurahansResponses
}

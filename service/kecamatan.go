package service

import (
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"gorm.io/gorm"
)

type KecamatanServiceInterface interface {
	FindKecamatan(requestId string) (kecamatansResponses []response.FindKecamatanResponse)
}

type KecamatanServiceImplementation struct {
	DB                           *gorm.DB
	Validate                     *validator.Validate
	Logger                       *logrus.Logger
	KecamatanRepositoryInterface repository.KecamatanRepositoryInterface
}

func NewKecamatanService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	kecamatanServiceInterface repository.KecamatanRepositoryInterface,
) KecamatanServiceInterface {
	return &KecamatanServiceImplementation{
		DB:                           db,
		Validate:                     validate,
		Logger:                       logger,
		KecamatanRepositoryInterface: kecamatanServiceInterface,
	}
}

func (service *KecamatanServiceImplementation) FindKecamatan(requestId string) (kecamatansResponses []response.FindKecamatanResponse) {
	kecamatans, err := service.KecamatanRepositoryInterface.FindKecamatan(service.DB)
	exceptions.PanicIfRecordNotFound(err, requestId, []string{"Data not found"}, service.Logger)
	kecamatansResponses = response.ToFindKecamatanResponse(kecamatans)
	return kecamatansResponses
}

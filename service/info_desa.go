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

type InfoDesaServiceInterface interface {
	FindInfoDesaByIdDesa(requestId string, idDesa string) (infoDesaResponses []response.FindInfoDesaByIdDesaResponse)
}

type InfoDesaServiceImplementation struct {
	DB                          *gorm.DB
	Validate                    *validator.Validate
	Logger                      *logrus.Logger
	InfoDesaRepositoryInterface repository.InfoDesaRepositoryInterface
}

func NewInfoDesaService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	infoDesaServiceInterface repository.InfoDesaRepositoryInterface,
) InfoDesaServiceInterface {
	return &InfoDesaServiceImplementation{
		DB:                          db,
		Validate:                    validate,
		Logger:                      logger,
		InfoDesaRepositoryInterface: infoDesaServiceInterface,
	}
}

func (service *InfoDesaServiceImplementation) FindInfoDesaByIdDesa(requestId string, idDesa string) (infoDesaResponses []response.FindInfoDesaByIdDesaResponse) {
	infoDesas, err := service.InfoDesaRepositoryInterface.FindInfoDesaByIdDesa(service.DB, idDesa)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(infoDesas) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("not found"), requestId, []string{"data not found"}, service.Logger)
	}
	infoDesaResponses = response.ToFindInfoDesaByIdDesaResponse(infoDesas)
	return infoDesaResponses
}

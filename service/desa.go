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

type DesaServiceInterface interface {
	FindDesaByIdKelu(requestId string, idKelu int) (desaResponses []response.FindDesaByIdKeluResponse)
}

type DesaServiceImplementation struct {
	DB                      *gorm.DB
	Validate                *validator.Validate
	Logger                  *logrus.Logger
	DesaRepositoryInterface repository.DesaRepositoryInterface
}

func NewDesaService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	desaServiceInterface repository.DesaRepositoryInterface,
) DesaServiceInterface {
	return &DesaServiceImplementation{
		DB:                      db,
		Validate:                validate,
		Logger:                  logger,
		DesaRepositoryInterface: desaServiceInterface,
	}
}

func (service *DesaServiceImplementation) FindDesaByIdKelu(requestId string, idKelu int) (desaResponses []response.FindDesaByIdKeluResponse) {
	desas, err := service.DesaRepositoryInterface.FindDesaByIdKelu(service.DB, idKelu)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(desas) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("not found"), requestId, []string{"data not found"}, service.Logger)
	}
	desaResponses = response.ToFindDesaByIdKeluResponse(desas)
	return desaResponses
}

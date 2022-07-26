package service

import (
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"gorm.io/gorm"
)

type PointServiceInterface interface {
	FindPointByUser(requestId string, idUser string) (pointResponse response.FindPointByUserResponse)
}

type PointServiceImplementation struct {
	DB                       *gorm.DB
	Validate                 *validator.Validate
	Logger                   *logrus.Logger
	PointRepositoryInterface repository.PointRepositoryInterface
}

func NewPointService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	pointServiceInterface repository.PointRepositoryInterface,
) PointServiceInterface {
	return &PointServiceImplementation{
		DB:                       db,
		Validate:                 validate,
		Logger:                   logger,
		PointRepositoryInterface: pointServiceInterface,
	}
}

func (service *PointServiceImplementation) FindPointByUser(requestId string, idUser string) (pointResponse response.FindPointByUserResponse) {
	point, err := service.PointRepositoryInterface.FindPointByUser(service.DB, idUser)
	exceptions.PanicIfRecordNotFound(err, requestId, []string{"Data not found"}, service.Logger)
	pointResponse = response.ToFindPointByUserResponse(point)
	return pointResponse
}

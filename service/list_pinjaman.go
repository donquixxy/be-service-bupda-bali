package service

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"gorm.io/gorm"
)

type ListPinjamanServiceInterface interface {
	FindListPinjamanByUser(requestId, IdUser string) (listPinjamanResponses []response.ListPinjamanResponse)
	FindListPinjamanByIdPinjaman(requestId, IdPinjaman string) (listPinjamanResponse response.ListPinjamanResponse)
}

type listPinjamanService struct {
	DB                              *gorm.DB
	Logger                          *logrus.Logger
	ListPinjamanRepositoryInterface repository.ListPinjamanRepositoryInterface
}

func NewListPinjamanService(
	db *gorm.DB,
	logger *logrus.Logger,
	listPinjamanServiceInterface repository.ListPinjamanRepositoryInterface,
) ListPinjamanServiceInterface {
	return &listPinjamanService{
		DB:                              db,
		Logger:                          logger,
		ListPinjamanRepositoryInterface: listPinjamanServiceInterface,
	}
}

func (service *listPinjamanService) FindListPinjamanByUser(requestId, IdUser string) (listPinjamanResponses []response.ListPinjamanResponse) {
	listPinjaman, err := service.ListPinjamanRepositoryInterface.FindListPinjamanByIdUser(service.DB, IdUser)
	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{"error when get list pinjaman by user " + err.Error()}, service.Logger)
	}
	if len(listPinjaman) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("list pinjaman not found"), requestId, []string{"list pinjaman not found"}, service.Logger)
	}
	listPinjamanResponses = response.ToListPinjamanResponses(listPinjaman)
	return listPinjamanResponses
}

func (service *listPinjamanService) FindListPinjamanByIdPinjaman(requestId, IdPinjaman string) (listPinjamanResponse response.ListPinjamanResponse) {
	listPinjaman, err := service.ListPinjamanRepositoryInterface.FindListPinjamanById(service.DB, IdPinjaman)
	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{"error when get list pinjaman by id pinjaman " + err.Error()}, service.Logger)
	}

	if len(listPinjaman.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("list pinjaman not found"), requestId, []string{"list pinjaman not found"}, service.Logger)
	}

	listPinjamanResponse = response.ToListPinjamanResponse(listPinjaman)
	return listPinjamanResponse
}

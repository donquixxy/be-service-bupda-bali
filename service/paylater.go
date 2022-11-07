package service

import (
	"errors"
	"log"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	invelirepository "github.com/tensuqiuwulu/be-service-bupda-bali/repository/inveli_repository"
	"github.com/tensuqiuwulu/be-service-bupda-bali/utilities"
	"gorm.io/gorm"
)

type PaylaterServiceInterface interface {
	CreatePaylater(requestId string, idUser string, requestCreatePaylater *request.CreatePaylaterRequest) error
}

type PaylaterServiceImplementation struct {
	DB                           *gorm.DB
	Validate                     *validator.Validate
	Logger                       *logrus.Logger
	UserRepositoryInterface      repository.UserRepositoryInterface
	InveliAPIRepositoryInterface invelirepository.InveliAPIRepositoryInterface
}

func NewPaylaterService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	userRepositoryInterface repository.UserRepositoryInterface,
	inveliAPIRepositoryInterface invelirepository.InveliAPIRepositoryInterface,
) PaylaterServiceInterface {
	return &PaylaterServiceImplementation{
		DB:                           db,
		Validate:                     validate,
		Logger:                       logger,
		UserRepositoryInterface:      userRepositoryInterface,
		InveliAPIRepositoryInterface: inveliAPIRepositoryInterface,
	}
}

func (service *PaylaterServiceImplementation) CreatePaylater(requestId string, idUser string, requestCreatePaylater *request.CreatePaylaterRequest) error {
	var err error
	var userAccount *entity.UserAccount

	request.ValidateRequest(service.Validate, requestCreatePaylater, requestId, service.Logger)

	user, err := service.UserRepositoryInterface.FindUserById2(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(user.Id) == 0 {
		exceptions.PanicIfBadRequest(errors.New("user not found"), requestId, []string{"user account not found"}, service.Logger)
	}

	userAccount, err = service.UserRepositoryInterface.GetUserAccountByID(service.DB, idUser)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(userAccount.Id) == 0 {

		accountInfo, _ := service.InveliAPIRepositoryInterface.GetAccountInfo(user.InveliAccessToken, user.InveliIDMember)
		if accountInfo == nil {
			exceptions.PanicIfBadRequest(errors.New("akun belum aktif"), requestId, []string{"akun belum aktif"}, service.Logger)
		}

		log.Println("accountInfo", accountInfo)

		var userAccounts []*entity.UserAccount
		for _, account := range accountInfo {
			userAccount := &entity.UserAccount{}
			userAccount.Id = utilities.RandomUUID()
			userAccount.IdUser = user.Id
			userAccount.IdAccount = account.ID
			userAccount.AccountName = account.AccountName2
			userAccount.IdProduct = account.ProductID
			userAccounts = append(userAccounts, userAccount)
		}

		errorr := service.UserRepositoryInterface.SaveUserAccount(service.DB, userAccounts)
		if errorr != nil {
			exceptions.PanicIfBadRequest(errors.New("gagal simpan user account"), "requestId", []string{"Failed Save User Account"}, service.Logger)
		}

	}
	userAccount, _ = service.UserRepositoryInterface.GetUserAccountByID(service.DB, idUser)

	service.InveliAPIRepositoryInterface.InveliCreatePaylater(user.InveliAccessToken, user.InveliIDMember, userAccount.IdAccount, requestCreatePaylater.Amount)

	// Fungsi untuk save history paylater ke database

	return nil

}

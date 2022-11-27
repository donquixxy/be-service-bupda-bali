package service

import (
	"fmt"
	"log"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	invelirepository "github.com/tensuqiuwulu/be-service-bupda-bali/repository/inveli_repository"
	"gorm.io/gorm"
)

type PaylaterServiceInterface interface {
	CreatePaylater(requestId string, idUser string, requestCreatePaylater *request.CreatePaylaterRequest) error
	GetTagihanPaylater(requestId string, idUser string) (tagihanPaylaterResponse []response.FindTagihanPaylater)
	GetOrderPaylaterPerBulan(requestId string, idUser string) (orderPaylaterPerBulanResponse []response.GetRiwayatPaylaterPerbulanResponse)
	GetOrderPaylaterByMonth(requestId string, idUser string, month int) (orderResponse []response.FindOrderByUserResponse)
	GetPembayaranTransaksiByIdUser(requestId, idUser, indexDate string) (response response.FindDetailPyamentPaylater)
}

type PaylaterServiceImplementation struct {
	DB                                *gorm.DB
	Validate                          *validator.Validate
	Logger                            *logrus.Logger
	UserRepositoryInterface           repository.UserRepositoryInterface
	InveliAPIRepositoryInterface      invelirepository.InveliAPIRepositoryInterface
	OrderRepositoryInterface          repository.OrderRepositoryInterface
	PaymentHistoryRepositoryInterface repository.PaymentHistoryRepositoryInterface
}

func NewPaylaterService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	userRepositoryInterface repository.UserRepositoryInterface,
	inveliAPIRepositoryInterface invelirepository.InveliAPIRepositoryInterface,
	orderRepositoryInterface repository.OrderRepositoryInterface,
	paymentHistoryRepositoryInterface repository.PaymentHistoryRepositoryInterface,
) PaylaterServiceInterface {
	return &PaylaterServiceImplementation{
		DB:                                db,
		Validate:                          validate,
		Logger:                            logger,
		UserRepositoryInterface:           userRepositoryInterface,
		InveliAPIRepositoryInterface:      inveliAPIRepositoryInterface,
		OrderRepositoryInterface:          orderRepositoryInterface,
		PaymentHistoryRepositoryInterface: paymentHistoryRepositoryInterface,
	}
}

func (service *PaylaterServiceImplementation) GetPembayaranTransaksiByIdUser(requestId, idUser, indexDate string) (response response.FindDetailPyamentPaylater) {
	paymentHistory, err := service.PaymentHistoryRepositoryInterface.FindPaymentHistoryById(service.DB, idUser, indexDate)
	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{"payment history not found"}, service.Logger)
	}

	response.Date = paymentHistory.TglPembayaran.Time
	response.NoTransaksi = paymentHistory.NoTransaksi
	response.Tagihan = paymentHistory.JmlTagihan
	response.BiayaAdmin = paymentHistory.BiayaAdmin
	response.Bunga = paymentHistory.BungaPinjaman
	response.Total = paymentHistory.Total
	return response
}

func (service *PaylaterServiceImplementation) GetOrderPaylaterByMonth(requestId string, idUser string, month int) (orderResponse []response.FindOrderByUserResponse) {
	order, err, _, _ := service.OrderRepositoryInterface.GetOrderPaylaterPerBulan(service.DB, idUser, month)
	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{"order not found"}, service.Logger)
	}

	orderResponse = response.ToFindOrderByUserResponse(order)
	return orderResponse
}

func (service *PaylaterServiceImplementation) GetOrderPaylaterPerBulan(requestId string, idUser string) (orderPaylaterPerBulanResponse []response.GetRiwayatPaylaterPerbulanResponse) {
	for i := 1; i <= 12; i++ {

		orderPaylaterPerBulan, _, start, end := service.OrderRepositoryInterface.GetOrderPaylaterPerBulan(service.DB, idUser, i)
		fmt.Println("data", orderPaylaterPerBulan)
		if len(orderPaylaterPerBulan) == 0 {
			continue
		}

		var responseData = response.GetRiwayatPaylaterPerbulanResponse{}
		for _, order := range orderPaylaterPerBulan {
			responseData.StartDate = start
			responseData.EndDate = end
			responseData.Month = i
			responseData.TotalBayar = responseData.TotalBayar + order.TotalBill
		}

		orderPaylaterPerBulanResponse = append(orderPaylaterPerBulanResponse, responseData)
	}

	return orderPaylaterPerBulanResponse
}

func (service *PaylaterServiceImplementation) GetTagihanPaylater(requestId string, idUser string) (tagihanPaylaterResponse []response.FindTagihanPaylater) {
	user, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{"user not found"}, service.Logger)
	}

	tagihanPaylater, err := service.InveliAPIRepositoryInterface.GetTagihanPaylater(user.User.InveliIDMember, user.User.InveliAccessToken)
	if err != nil {
		log.Println("error get tagihan inveli", err.Error())
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	tagihanPaylaterResponse = response.ToFindTagihanPaylater(tagihanPaylater)
	return tagihanPaylaterResponse
}

func (service *PaylaterServiceImplementation) CreatePaylater(requestId string, idUser string, requestCreatePaylater *request.CreatePaylaterRequest) error {
	// var err error
	// var userAccount *entity.UserAccount

	// request.ValidateRequest(service.Validate, requestCreatePaylater, requestId, service.Logger)

	// user, err := service.UserRepositoryInterface.FindUserById2(service.DB, idUser)
	// exceptions.PanicIfError(err, requestId, service.Logger)
	// if len(user.Id) == 0 {
	// 	exceptions.PanicIfBadRequest(errors.New("user not found"), requestId, []string{"user account not found"}, service.Logger)
	// }

	// userAccount, err = service.UserRepositoryInterface.GetUserAccountByID(service.DB, idUser)
	// exceptions.PanicIfError(err, requestId, service.Logger)
	// if len(userAccount.Id) == 0 {

	// 	accountInfo, _ := service.InveliAPIRepositoryInterface.GetAccountInfo(user.InveliAccessToken, user.InveliIDMember)
	// 	if accountInfo == nil {
	// 		exceptions.PanicIfBadRequest(errors.New("akun belum aktif"), requestId, []string{"akun belum aktif"}, service.Logger)
	// 	}

	// 	log.Println("accountInfo", accountInfo)

	// 	var userAccounts []*entity.UserAccount
	// 	for _, account := range accountInfo {
	// 		userAccount := &entity.UserAccount{}
	// 		userAccount.Id = utilities.RandomUUID()
	// 		userAccount.IdUser = user.Id
	// 		userAccount.IdAccount = account.ID
	// 		userAccount.AccountName = account.AccountName2
	// 		userAccount.IdProduct = account.ProductID
	// 		userAccounts = append(userAccounts, userAccount)
	// 	}

	// 	errorr := service.UserRepositoryInterface.SaveUserAccount(service.DB, userAccounts)
	// 	if errorr != nil {
	// 		exceptions.PanicIfBadRequest(errors.New("gagal simpan user account"), "requestId", []string{"Failed Save User Account"}, service.Logger)
	// 	}

	// }
	// userAccount, _ = service.UserRepositoryInterface.GetUserAccountByID(service.DB, idUser)

	// service.InveliAPIRepositoryInterface.InveliCreatePaylater(user.InveliAccessToken, user.InveliIDMember, userAccount.IdAccount, requestCreatePaylater.Amount)

	// // Fungsi untuk save history paylater ke database

	return nil

}

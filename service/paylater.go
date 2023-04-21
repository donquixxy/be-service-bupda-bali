package service

import (
	"log"
	"time"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	invelirepository "github.com/tensuqiuwulu/be-service-bupda-bali/repository/inveli_repository"
	"gorm.io/gorm"
)

type PaylaterServiceInterface interface {
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
	AuthServiceInterface              AuthServiceInterface
}

func NewPaylaterService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	userRepositoryInterface repository.UserRepositoryInterface,
	inveliAPIRepositoryInterface invelirepository.InveliAPIRepositoryInterface,
	orderRepositoryInterface repository.OrderRepositoryInterface,
	paymentHistoryRepositoryInterface repository.PaymentHistoryRepositoryInterface,
	authServiceInterface AuthServiceInterface,
) PaylaterServiceInterface {
	return &PaylaterServiceImplementation{
		DB:                                db,
		Validate:                          validate,
		Logger:                            logger,
		UserRepositoryInterface:           userRepositoryInterface,
		InveliAPIRepositoryInterface:      inveliAPIRepositoryInterface,
		OrderRepositoryInterface:          orderRepositoryInterface,
		PaymentHistoryRepositoryInterface: paymentHistoryRepositoryInterface,
		AuthServiceInterface:              authServiceInterface,
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

	order, _, _, err := service.OrderRepositoryInterface.GetOrderPaylaterPerBulan(service.DB, idUser, month)
	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{"order not found"}, service.Logger)
	}

	orderResponse = response.ToFindOrderByUserResponse(order)
	return orderResponse
}

func (service *PaylaterServiceImplementation) GetOrderPaylaterPerBulan(requestId string, idUser string) (orderPaylaterPerBulanResponse []response.GetRiwayatPaylaterPerbulanResponse) {

	yearNow := time.Now().Year()
	monthNow := time.Now().Month()

	for i := 12; i > 0; i-- {
		var responseData = response.GetRiwayatPaylaterPerbulanResponse{}
		orderPaylaterPerBulan, _ := service.OrderRepositoryInterface.FindOrderTotalPaylaterByMonth(service.DB, idUser, int(monthNow))
		responseData.StartDate = time.Date(yearNow, monthNow, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02 15:04:05")
		responseData.EndDate = time.Date(yearNow, monthNow+1, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02 15:04:05")
		responseData.Month = int(monthNow)

		for _, order := range orderPaylaterPerBulan {
			responseData.TotalBayar = responseData.TotalBayar + order.TotalBill
		}

		orderPaylaterPerBulanResponse = append(orderPaylaterPerBulanResponse, responseData)

		monthNow = monthNow - 1
		if monthNow == 0 {
			monthNow = 12
			yearNow = yearNow - 1
		}

		if responseData.TotalBayar == 0 {
			orderPaylaterPerBulanResponse = orderPaylaterPerBulanResponse[:len(orderPaylaterPerBulanResponse)-1]
		}

	}

	return orderPaylaterPerBulanResponse
}

func (service *PaylaterServiceImplementation) GetTagihanPaylater(requestId string, idUser string) (tagihanPaylaterResponse []response.FindTagihanPaylater) {
	// user, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	// if err != nil {
	// 	exceptions.PanicIfBadRequest(err, requestId, []string{"user not found"}, service.Logger)
	// }

	// _, err = service.InveliAPIRepositoryInterface.GetLastLoanIdPaylater(user.User.InveliIDMember, user.User.InveliAccessToken)
	// if err != nil {
	// 	log.Println("error get tagihan inveli", err.Error())
	// 	exceptions.PanicIfError(err, requestId, service.Logger)
	// }

	// log.Println("tagihanPaylater = ", tagihanPaylater[0].LoanAccountID)

	// if tagihanPaylater == nil {
	// 	if user.User.StatusPaylater == 2 {
	// 		go service.AuthServiceInterface.FirstTimeLoginInveli(user.User.Phone, user.User.InveliPassword)
	// 	}
	// }

	tagihanPaylater, err := service.OrderRepositoryInterface.FindOrderPaylaterUnpaidById(service.DB, idUser)
	if err != nil {
		log.Println("error get tagihan inveli", err.Error())
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	tagihanPaylaterResponse = response.ToFindTagihanPaylater(tagihanPaylater)
	return tagihanPaylaterResponse

}

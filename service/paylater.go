package service

import (
	"errors"
	"log"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/inveli"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	invelirepository "github.com/tensuqiuwulu/be-service-bupda-bali/repository/inveli_repository"
	"gorm.io/gorm"
)

type PaylaterServiceInterface interface {
	GetTagihanPaylater(requestId string, idUser string) (tagihanPaylaterResponse response.TotalTagihan)
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
	order, _, _, err := service.OrderRepositoryInterface.GetOrderPaylaterPerBulan(service.DB, idUser, month)
	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{"order not found"}, service.Logger)
	}

	orderResponse = response.ToFindOrderByUserResponse(order)
	return orderResponse
}

func (service *PaylaterServiceImplementation) GetOrderPaylaterPerBulan(requestId string, idUser string) (orderPaylaterPerBulanResponse []response.GetRiwayatPaylaterPerbulanResponse) {
	for i := 1; i <= 12; i++ {

		orderPaylaterPerBulan, start, end, _ := service.OrderRepositoryInterface.GetOrderPaylaterPerBulan(service.DB, idUser, i)
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

func (service *PaylaterServiceImplementation) GetTagihanPaylater(requestId string, idUser string) (tagihanPaylaterResponse response.TotalTagihan) {
	// log.Println("masuk ke tagihan paylater")
	user, err := service.UserRepositoryInterface.FindUserById(service.DB, idUser)
	if err != nil {
		exceptions.PanicIfBadRequest(err, requestId, []string{"user not found"}, service.Logger)
	}

	tagihanPaylater, err := service.InveliAPIRepositoryInterface.GetTagihanPaylater(user.User.InveliIDMember, user.User.InveliAccessToken)
	if err != nil {
		log.Println("error get tagihan inveli", err.Error())
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	// log.Println("tagihan", tagihanPaylater)

	count := 0
	for _, tagihan := range tagihanPaylater {
		if tagihan.IsPaid {
			continue
		}
		count++
	}

	// log.Println("count", count)

	if count == 0 {
		// log.Println("MASUK")
		// get riwayat pinjaman
		loanID, err := service.InveliAPIRepositoryInterface.GetRiwayatPinjaman(user.User.InveliAccessToken, user.User.InveliIDMember)
		if err != nil {
			log.Println("error get riwayat pinjaman", err.Error())
			exceptions.PanicIfError(err, requestId, service.Logger)
		}

		log.Println("loanID", loanID)

		if len(loanID) == 0 {
			exceptions.PanicIfBadRequest(errors.New("riwayat paylater not found"), requestId, []string{"riwayat paylater not found"}, service.Logger)
		}

		// log.Println("loanID", loanID)

		tunggakans := []inveli.TunggakanPaylater{}
		for _, id := range loanID {
			tunggakan, err := service.InveliAPIRepositoryInterface.GetTunggakan(id, user.User.InveliAccessToken)
			if err != nil {
				log.Println("error get tunggakan", err.Error())
				exceptions.PanicIfError(err, requestId, service.Logger)
			}

			if len(tunggakan) == 0 {
				exceptions.PanicIfBadRequest(errors.New("tunggakan not found"), requestId, []string{"tunggakan not found"}, service.Logger)
			}

			tunggakans = append(tunggakans, tunggakan...)
		}

		// log.Println("tunggakan", tunggakans)
		tagihanPaylaterResponse = response.ToFindTunggakanPaylater(tunggakans)
		return tagihanPaylaterResponse
	} else {
		tagihanPaylaterResponse = response.ToFindTagihanPaylater(tagihanPaylater)
		return tagihanPaylaterResponse
	}

}

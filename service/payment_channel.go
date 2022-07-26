package service

import (
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"gorm.io/gorm"
)

type PaymentChannelServiceInterface interface {
	FindPaymentChannel(requestId string) (paymentChannelResponses []response.FindPaymentChannelResponse)
}

type PaymentChannelServiceImplementation struct {
	DB                                *gorm.DB
	Validate                          *validator.Validate
	Logger                            *logrus.Logger
	PaymentChannelRepositoryInterface repository.PaymentChannelRepositoryInterface
}

func NewPaymentChannelService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	paymentChannelRepositoryInterface repository.PaymentChannelRepositoryInterface,
) PaymentChannelServiceInterface {
	return &PaymentChannelServiceImplementation{
		DB:                                db,
		Validate:                          validate,
		Logger:                            logger,
		PaymentChannelRepositoryInterface: paymentChannelRepositoryInterface,
	}
}

func (service *PaymentChannelServiceImplementation) FindPaymentChannel(requestId string) (paymentChannelResponses []response.FindPaymentChannelResponse) {
	paymentChannelResponse, err := service.PaymentChannelRepositoryInterface.FindPaymentChannel(service.DB)
	exceptions.PanicIfRecordNotFound(err, requestId, []string{"Data not found"}, service.Logger)
	paymentChannelResponses = response.ToFindPaymentChannelResponse(paymentChannelResponse)
	return paymentChannelResponses
}

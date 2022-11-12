package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type PaylaterControllerInterface interface {
	CreatePaylater(c echo.Context) error
	PayPaylater(c echo.Context) error
	GetTagihanPaylater(c echo.Context) error
}

type PaylaterControllerImplementation struct {
	logger                   *logrus.Logger
	PaylaterServiceInterface service.PaylaterServiceInterface
	PaymentServiceInterface  service.PaymentServiceInterface
}

func NewPaylaterController(
	logger *logrus.Logger,
	paylaterServiceInterface service.PaylaterServiceInterface,
	paymentServiceInterface service.PaymentServiceInterface,
) PaylaterControllerInterface {
	return &PaylaterControllerImplementation{
		logger:                   logger,
		PaylaterServiceInterface: paylaterServiceInterface,
		PaymentServiceInterface:  paymentServiceInterface,
	}
}

func (controller *PaylaterControllerImplementation) GetTagihanPaylater(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdUser := middleware.TokenClaimsIdUser(c)
	tagihanResponse := controller.PaylaterServiceInterface.GetTagihanPaylater(requestId, IdUser)
	if condition := len(tagihanResponse) == 0; condition {
		exceptions.PanicIfRecordNotFound(errors.New("tagihan paylater not found"), requestId, []string{"tagihan paylater not found"}, controller.logger)
	}

	responses := response.Response{Code: 200, Mssg: "success", Data: tagihanResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PaylaterControllerImplementation) PayPaylater(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdUser := middleware.TokenClaimsIdUser(c)
	controller.PaymentServiceInterface.PayPaylater(requestId, IdUser)
	responses := response.Response{Code: 201, Mssg: "success", Data: "success pay paylater", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PaylaterControllerImplementation) CreatePaylater(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdUser := middleware.TokenClaimsIdUser(c)
	request := request.ReadFromCreatePaylaterRequestBody(c, requestId, controller.logger)
	controller.PaylaterServiceInterface.CreatePaylater(requestId, IdUser, request)
	responses := response.Response{Code: 201, Mssg: "success", Data: "success create paylater", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

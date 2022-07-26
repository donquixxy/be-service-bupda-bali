package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type OtpManagerControllerInterface interface {
	SendOtpBySms(c echo.Context) error
	VerifyOtp(c echo.Context) error
}

type OtpManagerControllerImplementation struct {
	ConfigurationWebserver     config.Webserver
	Logger                     *logrus.Logger
	OtpManagerServiceInterface service.OtpManagerServiceInterface
}

func NewOtpManagerController(
	logger *logrus.Logger,
	otpManagerServiceInterface service.OtpManagerServiceInterface,
) OtpManagerControllerInterface {
	return &OtpManagerControllerImplementation{
		Logger:                     logger,
		OtpManagerServiceInterface: otpManagerServiceInterface,
	}
}

func (controller *OtpManagerControllerImplementation) SendOtpBySms(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	request := request.ReadFromSendOtpBySmsRequestBody(c, requestId, controller.Logger)
	controller.OtpManagerServiceInterface.SendOtpBySms(requestId, request)
	respon := response.Response{Code: 200, Mssg: "success", Data: nil, Error: []string{}}
	return c.JSON(http.StatusOK, respon)
}

func (controller *OtpManagerControllerImplementation) VerifyOtp(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	request := request.ReadFromVerifyOtpRequestBody(c, requestId, controller.Logger)
	token := controller.OtpManagerServiceInterface.VerifyOtp(requestId, request)
	fmt.Println("masuk", token)
	respon := response.Response{Code: 200, Mssg: "success", Data: token, Error: []string{}}
	return c.JSON(http.StatusOK, respon)
}

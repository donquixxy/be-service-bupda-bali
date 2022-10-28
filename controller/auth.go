package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type AuthControllerInterface interface {
	Login(c echo.Context) error
	FirstTimeLoginInveli(c echo.Context) error
	FirstTimeUbahPasswordInveli(c echo.Context) error
	NewToken(c echo.Context) error
}

type AuthControllerImplementation struct {
	ConfigurationWebserver config.Webserver
	Logger                 *logrus.Logger
	AuthServiceInterface   service.AuthServiceInterface
}

func NewAuthController(
	logger *logrus.Logger,
	authServiceInterface service.AuthServiceInterface,
) AuthControllerInterface {
	return &AuthControllerImplementation{
		Logger:               logger,
		AuthServiceInterface: authServiceInterface,
	}
}

func (controller *AuthControllerImplementation) Login(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	request := request.ReadFromLoginRequestBody(c, requestId, controller.Logger)
	loginResponse := controller.AuthServiceInterface.Login(requestId, request)
	respon := response.Response{Code: 200, Mssg: "success", Data: loginResponse, Error: []string{}}
	return c.JSON(http.StatusOK, respon)
}

func (controller *AuthControllerImplementation) FirstTimeLoginInveli(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	request := request.ReadFromLoginInveliRequestBody(c, requestId, controller.Logger)
	loginInveliResponse := controller.AuthServiceInterface.FirstTimeLoginInveli(requestId, request)
	respon := response.Response{Code: 200, Mssg: "success", Data: loginInveliResponse, Error: []string{}}
	return c.JSON(http.StatusOK, respon)
}

func (controller *AuthControllerImplementation) FirstTimeUbahPasswordInveli(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	request := request.ReadFromUpdateUserPasswordInveliRequestBody(c, requestId, controller.Logger)
	controller.AuthServiceInterface.FirstTimeUbahPasswordInveli(requestId, request)
	respon := response.Response{Code: 200, Mssg: "success", Data: nil, Error: []string{}}
	return c.JSON(http.StatusOK, respon)
}

func (controller *AuthControllerImplementation) NewToken(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	refreshToken := c.FormValue("refresh_token")
	token := controller.AuthServiceInterface.NewToken(requestId, refreshToken)
	respon := response.Response{Code: 200, Mssg: "success", Data: token, Error: []string{}}
	return c.JSON(http.StatusOK, respon)
}

package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type UserShippingAddressControllerInterface interface {
	FindUserShippingAddress(c echo.Context) error
	CreateUserShippingAddress(c echo.Context) error
	DeleteUserShippingAddress(c echo.Context) error
}

type UserShippingAddressControllerImplementation struct {
	Logger                              *logrus.Logger
	UserShippingAddressServiceInterface service.UserShippingAddressServiceInterface
}

func NewUserShippingAddressController(
	logger *logrus.Logger,
	userShippingAddressServiceInterface service.UserShippingAddressServiceInterface,
) UserShippingAddressControllerInterface {
	return &UserShippingAddressControllerImplementation{
		Logger:                              logger,
		UserShippingAddressServiceInterface: userShippingAddressServiceInterface,
	}
}

func (controller *UserShippingAddressControllerImplementation) DeleteUserShippingAddress(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	userShippingAddressRequest := request.ReadFromUserShippingAddressRequestBody(c, requestId, controller.Logger)
	controller.UserShippingAddressServiceInterface.DeleteUserShippingAddress(requestId, userShippingAddressRequest)
	responses := response.Response{Code: 200, Mssg: "success", Data: nil, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserShippingAddressControllerImplementation) CreateUserShippingAddress(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	request := request.ReadFromCreateUserShippingAddressRequestBody(c, requestId, controller.Logger)
	controller.UserShippingAddressServiceInterface.CreateUserShippingAddress(requestId, idUser, request)
	responses := response.Response{Code: 200, Mssg: "success", Data: nil, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserShippingAddressControllerImplementation) FindUserShippingAddress(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	userAddressResponses := controller.UserShippingAddressServiceInterface.FindUserShippingAddressByIdUser(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: userAddressResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

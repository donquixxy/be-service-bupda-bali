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

type UserControllerInterface interface {
	CreateUserNonSuveyed(c echo.Context) error
	FindUserById(c echo.Context) error
	DeleteUserById(c echo.Context) error
}

type UserControllerImplementation struct {
	Logger               *logrus.Logger
	UserServiceInterface service.UserServiceInterface
}

func NewUserController(
	logger *logrus.Logger,
	userServiceInterface service.UserServiceInterface,
) UserControllerInterface {
	return &UserControllerImplementation{
		UserServiceInterface: userServiceInterface,
	}
}

func (controller *UserControllerImplementation) CreateUserNonSuveyed(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	request := request.ReadFromCreateUserRequestBody(c, requestId, controller.Logger)
	controller.UserServiceInterface.CreateUserNonSuveyed(requestId, request)
	responses := response.Response{Code: 200, Mssg: "success", Data: nil, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) FindUserById(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	userResponse := controller.UserServiceInterface.FindUserById(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: userResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) DeleteUserById(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	controller.UserServiceInterface.DeleteUserById(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: "delete user success", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

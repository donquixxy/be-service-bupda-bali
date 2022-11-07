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

type PaylaterControllerInterface interface {
	CreatePaylater(c echo.Context) error
}

type PaylaterControllerImplementation struct {
	logger                   *logrus.Logger
	PaylaterServiceInterface service.PaylaterServiceInterface
}

func NewPaylaterController(
	logger *logrus.Logger,
	paylaterServiceInterface service.PaylaterServiceInterface,
) PaylaterControllerInterface {
	return &PaylaterControllerImplementation{
		logger:                   logger,
		PaylaterServiceInterface: paylaterServiceInterface,
	}
}

func (controller *PaylaterControllerImplementation) CreatePaylater(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdUser := middleware.TokenClaimsIdUser(c)
	request := request.ReadFromCreatePaylaterRequestBody(c, requestId, controller.logger)
	controller.PaylaterServiceInterface.CreatePaylater(requestId, IdUser, request)
	responses := response.Response{Code: 201, Mssg: "success", Data: "success create paylater", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

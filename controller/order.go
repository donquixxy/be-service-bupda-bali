package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type OrderControllerInterface interface {
	CreateOrder(c echo.Context) error
	FindOrderByUser(c echo.Context) error
	FindOrderById(c echo.Context) error
	CancelOrderById(c echo.Context) error
	CompleteOrderById(c echo.Context) error
}

type OrderControllerImplementation struct {
	Logger                *logrus.Logger
	OrderServiceInterface service.OrderServiceInterface
}

func NewOrderController(
	logger *logrus.Logger,
	orderServiceInterface service.OrderServiceInterface,
) OrderControllerInterface {
	return &OrderControllerImplementation{
		OrderServiceInterface: orderServiceInterface,
	}
}

func (controller *OrderControllerImplementation) CreateOrder(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	idDesa := middleware.TokenClaimsIdDesa(c)
	accountType := middleware.TokenClaimsAccountType(c)
	request := request.ReadFromCreateOrderRequestBody(c, requestId, controller.Logger)
	orderResponse := controller.OrderServiceInterface.CreateOrder(requestId, idUser, idDesa, accountType, request)
	responses := response.Response{Code: 201, Mssg: "success", Data: orderResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *OrderControllerImplementation) FindOrderByUser(c echo.Context) error {
	var orderStatus int
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	paramOrder := c.QueryParam("order_status")
	if paramOrder == "" {
		orderStatus = -1
	} else {
		orderStatus, _ = strconv.Atoi(paramOrder)
	}
	orderResponses := controller.OrderServiceInterface.FindOrderByUser(requestId, idUser, orderStatus)
	responses := response.Response{Code: 200, Mssg: "success", Data: orderResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *OrderControllerImplementation) FindOrderById(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idOrder := c.QueryParam("id_order")
	orderResponse := controller.OrderServiceInterface.FindOrderById(requestId, idOrder)
	responses := response.Response{Code: 200, Mssg: "success", Data: orderResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *OrderControllerImplementation) CancelOrderById(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	orderRequest := request.ReadFromOrderIdRequestBody(c, requestId, controller.Logger)
	controller.OrderServiceInterface.CancelOrderById(requestId, orderRequest)
	response := response.Response{Code: 201, Mssg: "success", Data: nil, Error: []string{}}
	return c.JSON(http.StatusOK, response)
}

func (controller *OrderControllerImplementation) CompleteOrderById(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	orderRequest := request.ReadFromOrderIdRequestBody(c, requestId, controller.Logger)
	controller.OrderServiceInterface.CompleteOrderById(requestId, orderRequest)
	response := response.Response{Code: 201, Mssg: "success", Data: nil, Error: []string{}}
	return c.JSON(http.StatusOK, response)
}

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
	UpdateOrderPaymentStatus(c echo.Context) error
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
	var orderResponse interface{}
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	idDesa := middleware.TokenClaimsIdDesa(c)
	accountType := middleware.TokenClaimsAccountType(c)
	orderType, _ := strconv.Atoi(c.QueryParam("order_type"))
	productType := c.QueryParam("product_type")

	if orderType == 1 {
		request := request.ReadFromCreateOrderRequestBody(c, requestId, controller.Logger)
		orderResponse = controller.OrderServiceInterface.CreateOrderSembako(requestId, idUser, idDesa, accountType, request)
	} else if orderType == 2 {
		switch productType {
		case "prepaid_pulsa", "prepaid_data", "prepaid_pln":
			request := request.ReadFromCreateOrderPrepaidRequestBody(c, requestId, controller.Logger)
			orderResponse = controller.OrderServiceInterface.CreateOrderPrepaid(requestId, idUser, idDesa, productType, request)
		case "postpaid_pln":
			// request := request.ReadFromCreateOrderPostpaidPlnRequestBody(c, requestId, controller.Logger)
			// orderResponse = controller.OrderServiceInterface.CreateOrderPostpaidPln(requestId, idUser, idDesa, request)
		case "postpaid_pdam":
		}
	}
	responses := response.Response{Code: 201, Mssg: "success", Data: orderResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *OrderControllerImplementation) FindOrderByUser(c echo.Context) error {
	var orderStatus int
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	paramOrder := c.QueryParam("order_status")
	if paramOrder == "" {
		// -1 untuk menampilkan semua order
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

func (controller *OrderControllerImplementation) UpdateOrderPaymentStatus(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	orderRequest := request.ReadFromUpdatePaymentStatusOrderRequestBody(c, requestId, controller.Logger)
	controller.OrderServiceInterface.UpdatePaymentStatusOrder(requestId, orderRequest)
	response := response.Response{Code: 201, Mssg: "success", Data: nil, Error: []string{}}
	return c.JSON(http.StatusOK, response)
}
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

type CartControllerInterface interface {
	CreateCart(c echo.Context) error
	UpdateCart(c echo.Context) error
	FindCartByUser(c echo.Context) error
}

type CartControllerImplementation struct {
	Logger               *logrus.Logger
	CartServiceInterface service.CartServiceInterface
}

func NewCartController(
	logger *logrus.Logger,
	cartServiceInterface service.CartServiceInterface,
) CartControllerInterface {
	return &CartControllerImplementation{
		CartServiceInterface: cartServiceInterface,
	}
}

func (controller *CartControllerImplementation) CreateCart(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	request := request.ReadFromCreateCartRequestBody(c, requestId, controller.Logger)
	idCart := controller.CartServiceInterface.CreateCart(requestId, idUser, request)
	data := make(map[string]interface{})
	data["id_cart"] = idCart
	responses := response.Response{Code: 200, Mssg: "success", Data: data, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *CartControllerImplementation) UpdateCart(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	request := request.ReadFromUpdateCartRequestBody(c, requestId, controller.Logger)
	idCart := controller.CartServiceInterface.UpdateCart(requestId, request)
	data := make(map[string]interface{})
	data["id_cart"] = idCart
	responses := response.Response{Code: 200, Mssg: "success", Data: data, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *CartControllerImplementation) FindCartByUser(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	idDesa := middleware.TokenClaimsIdDesa(c)
	accountType := middleware.TokenClaimsAccountType(c)
	cartResponse := controller.CartServiceInterface.FindCartByUser(requestId, idUser, accountType, idDesa)
	responses := response.Response{Code: 200, Mssg: "success", Data: cartResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

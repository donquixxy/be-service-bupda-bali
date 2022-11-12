package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type PaymentChannelControllerInterface interface {
	FindPaymentChannel(c echo.Context) error
}

type PaymentChannelControllerImplementation struct {
	PaymentChannelServiceInterface service.PaymentChannelServiceInterface
}

func NewPaymentChannelController(
	paymentChannelServiceInterface service.PaymentChannelServiceInterface) PaymentChannelControllerInterface {
	return &PaymentChannelControllerImplementation{
		PaymentChannelServiceInterface: paymentChannelServiceInterface,
	}
}

func (controller *PaymentChannelControllerImplementation) FindPaymentChannel(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	paymentChannelResponses := controller.PaymentChannelServiceInterface.FindPaymentChannel(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: paymentChannelResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type PpobControllerInterface interface {
	GetPulsaPriceList(c echo.Context) error
}

type PpobControllerImplementation struct {
	PpobServiceInterface service.PpobServiceInterface
}

func NewPpobController(
	ppobServiceInterface service.PpobServiceInterface) PpobControllerInterface {
	return &PpobControllerImplementation{
		PpobServiceInterface: ppobServiceInterface,
	}
}

func (controller *PpobControllerImplementation) GetPulsaPriceList(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	numberPhone := c.QueryParam("phone")
	pulsaPriceListResponses := controller.PpobServiceInterface.GetPrepaidPulsaPriceList(requestId, numberPhone)
	responses := response.Response{Code: 200, Mssg: "success", Data: pulsaPriceListResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

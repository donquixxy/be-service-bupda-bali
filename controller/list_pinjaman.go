package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type ListPinjamanControllerInterface interface {
	FindListPinjamanByUser(c echo.Context) error
	FindListPinjamanById(c echo.Context) error
}

type ListPinjamanControllerImplementation struct {
	ListPinjamanServiceInterface service.ListPinjamanServiceInterface
}

func NewListPinjamanController(
	listPinjamanServiceInterface service.ListPinjamanServiceInterface,
) ListPinjamanControllerInterface {
	return &ListPinjamanControllerImplementation{
		ListPinjamanServiceInterface: listPinjamanServiceInterface,
	}
}

func (controller *ListPinjamanControllerImplementation) FindListPinjamanByUser(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	listPinjamanResponses := controller.ListPinjamanServiceInterface.FindListPinjamanByUser(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: listPinjamanResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *ListPinjamanControllerImplementation) FindListPinjamanById(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idPinjaman := c.QueryParam("id_pinjaman")
	listPinjamanResponse := controller.ListPinjamanServiceInterface.FindListPinjamanByIdPinjaman(requestId, idPinjaman)
	responses := response.Response{Code: 200, Mssg: "success", Data: listPinjamanResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

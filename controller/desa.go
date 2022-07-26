package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type DesaControllerInterface interface {
	FindDesaByIdKelu(c echo.Context) error
}

type DesaControllerImplementation struct {
	DesaServiceInterface service.DesaServiceInterface
}

func NewDesaController(
	desaServiceInterface service.DesaServiceInterface) DesaControllerInterface {
	return &DesaControllerImplementation{
		DesaServiceInterface: desaServiceInterface,
	}
}

func (controller *DesaControllerImplementation) FindDesaByIdKelu(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idKelu, _ := strconv.Atoi(c.QueryParam("id_kelurahan"))
	desaResponses := controller.DesaServiceInterface.FindDesaByIdKelu(requestId, idKelu)
	responses := response.Response{Code: 200, Mssg: "success", Data: desaResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

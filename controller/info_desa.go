package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type InfoDesaControllerInterface interface {
	FindInfoDesaByIdDesa(c echo.Context) error
}

type InfoDesaControllerImplementation struct {
	InfoDesaServiceInterface service.InfoDesaServiceInterface
}

func NewInfoDesaController(
	infoDesaServiceInterface service.InfoDesaServiceInterface) InfoDesaControllerInterface {
	return &InfoDesaControllerImplementation{
		InfoDesaServiceInterface: infoDesaServiceInterface,
	}
}

func (controller *InfoDesaControllerImplementation) FindInfoDesaByIdDesa(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idDesa := middleware.TokenClaimsIdDesa(c)
	infoDesaResponses := controller.InfoDesaServiceInterface.FindInfoDesaByIdDesa(requestId, idDesa)
	responses := response.Response{Code: 200, Mssg: "success", Data: infoDesaResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type PointControllerInterface interface {
	FindPointByUser(c echo.Context) error
}

type PointControllerImplementation struct {
	PointServiceInterface service.PointServiceInterface
}

func NewPointController(
	pointServiceInterface service.PointServiceInterface) PointControllerInterface {
	return &PointControllerImplementation{
		PointServiceInterface: pointServiceInterface,
	}
}

func (controller *PointControllerImplementation) FindPointByUser(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdDesa(c)
	pointResponse := controller.PointServiceInterface.FindPointByUser(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: pointResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

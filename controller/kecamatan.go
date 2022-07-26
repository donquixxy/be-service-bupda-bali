package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type KecamatanControllerInterface interface {
	FindKecamatan(c echo.Context) error
}

type KecamatanControllerImplementation struct {
	KecamatanServiceInterface service.KecamatanServiceInterface
}

func NewKecamatanController(
	kecamatanServiceInterface service.KecamatanServiceInterface) KecamatanControllerInterface {
	return &KecamatanControllerImplementation{
		KecamatanServiceInterface: kecamatanServiceInterface,
	}
}

func (controller *KecamatanControllerImplementation) FindKecamatan(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	kecamatanResponses := controller.KecamatanServiceInterface.FindKecamatan(requestId)
	responses := response.Response{Code: 200, Mssg: "success", Data: kecamatanResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

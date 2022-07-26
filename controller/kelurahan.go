package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type KelurahanControllerInterface interface {
	FindKelurahanByIdKeca(c echo.Context) error
}

type KelurahanControllerImplementation struct {
	KelurahanServiceInterface service.KelurahanServiceInterface
}

func NewKelurahanController(
	kelurahanServiceInterface service.KelurahanServiceInterface) KelurahanControllerInterface {
	return &KelurahanControllerImplementation{
		KelurahanServiceInterface: kelurahanServiceInterface,
	}
}

func (controller *KelurahanControllerImplementation) FindKelurahanByIdKeca(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idKeca, _ := strconv.Atoi(c.QueryParam("id_kecamatan"))
	kecamatanResponses := controller.KelurahanServiceInterface.FindKelurahanByIdKeca(requestId, idKeca)
	responses := response.Response{Code: 200, Mssg: "success", Data: kecamatanResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

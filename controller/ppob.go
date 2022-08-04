package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type PpobControllerInterface interface {
	GetPrepaidPulsaPriceList(c echo.Context) error
	GetPrepaidDataPriceList(c echo.Context) error
	GetPrepaidPlnPriceList(c echo.Context) error
	InquiryPrepaidPln(c echo.Context) error
	InquiryPostpaidPln(c echo.Context) error
	GetPostpaidPdamProduct(c echo.Context) error
	InquiryPostpaidPdam(c echo.Context) error
}

type PpobControllerImplementation struct {
	Logger               *logrus.Logger
	PpobServiceInterface service.PpobServiceInterface
}

func NewPpobController(
	logger *logrus.Logger,
	ppobServiceInterface service.PpobServiceInterface) PpobControllerInterface {
	return &PpobControllerImplementation{
		Logger:               logger,
		PpobServiceInterface: ppobServiceInterface,
	}
}

func (controller *PpobControllerImplementation) GetPrepaidPulsaPriceList(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	numberPhone := c.QueryParam("phone")
	pulsaPriceListResponses := controller.PpobServiceInterface.GetPrepaidPulsaPriceList(requestId, numberPhone)
	responses := response.Response{Code: 200, Mssg: "success", Data: pulsaPriceListResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PpobControllerImplementation) GetPrepaidDataPriceList(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	numberPhone := c.QueryParam("phone")
	pulsaPriceListResponses := controller.PpobServiceInterface.GetPrepaidDataPriceList(requestId, numberPhone)
	responses := response.Response{Code: 200, Mssg: "success", Data: pulsaPriceListResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PpobControllerImplementation) GetPrepaidPlnPriceList(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idPelanggan := c.QueryParam("id_pelanggan")
	plnPriceListResponses := controller.PpobServiceInterface.GetPrepaidPlnPriceList(requestId, idPelanggan)
	responses := response.Response{Code: 200, Mssg: "success", Data: plnPriceListResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PpobControllerImplementation) InquiryPrepaidPln(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	InquiryPrepaidPlnRequest := request.ReadFromInquiryPrepaidPlnRequestBody(c, requestId, controller.Logger)
	detailResponses := controller.PpobServiceInterface.InquiryPrepaidPln(requestId, InquiryPrepaidPlnRequest)
	responses := response.Response{Code: 200, Mssg: "success", Data: detailResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PpobControllerImplementation) InquiryPostpaidPln(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	InquiryPostpaidPlnRequest := request.ReadFromInquiryPostpaidPlnRequestBody(c, requestId, controller.Logger)
	detailResponses := controller.PpobServiceInterface.InquiryPostpaidPln(requestId, InquiryPostpaidPlnRequest)
	responses := response.Response{Code: 200, Mssg: "success", Data: detailResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PpobControllerImplementation) GetPostpaidPdamProduct(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	detailResponses := controller.PpobServiceInterface.GetPostpaidPdamProduct(requestId)
	responses := response.Response{Code: 200, Mssg: "success", Data: detailResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PpobControllerImplementation) InquiryPostpaidPdam(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	inquiryPostpaidPdamRequest := request.ReadFromInquiryPostpaidPdamRequestBody(c, requestId, controller.Logger)
	detailResponses := controller.PpobServiceInterface.InquiryPostpaidPdam(requestId, inquiryPostpaidPdamRequest)
	responses := response.Response{Code: 200, Mssg: "success", Data: detailResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

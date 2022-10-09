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

type MerchantControllerInterface interface {
	CreateMerchantApproveList(c echo.Context) error
	FindMerchantStatusApproveByUserResponse(c echo.Context) error
}

type MerchantControllerImplementation struct {
	Logger                   *logrus.Logger
	MerchantServiceInterface service.MerchantServiceInterface
}

func NewMerchantController(
	logger *logrus.Logger,
	merchantServiceInterface service.MerchantServiceInterface,
) MerchantControllerInterface {
	return &MerchantControllerImplementation{
		MerchantServiceInterface: merchantServiceInterface,
	}
}

func (controller *MerchantControllerImplementation) CreateMerchantApproveList(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	idDesa := middleware.TokenClaimsIdDesa(c)
	request := request.ReadFromCreateMerchantApproveListRequestBody(c, requestId, controller.Logger)
	controller.MerchantServiceInterface.CreateMerchantApproveList(requestId, idUser, idDesa, request)
	responses := response.Response{Code: 200, Mssg: "success", Data: "Create Merchant Approve List Success", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *MerchantControllerImplementation) FindMerchantStatusApproveByUserResponse(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	merchantResponse := controller.MerchantServiceInterface.FindMerchantStatusApproveByUser(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: merchantResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

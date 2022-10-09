package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type BannerControllerInterface interface {
	FindBannerByDesa(c echo.Context) error
	FindBannerAll(c echo.Context) error
}

type BannerControllerImplementation struct {
	Logger                 *logrus.Logger
	BannerServiceInterface service.BannerServiceInterface
}

func NewBannerController(
	logger *logrus.Logger,
	bannerServiceInterface service.BannerServiceInterface,
) BannerControllerInterface {
	return &BannerControllerImplementation{
		Logger:                 logger,
		BannerServiceInterface: bannerServiceInterface,
	}
}

func (controller *BannerControllerImplementation) FindBannerAll(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	bannerResponses := controller.BannerServiceInterface.FindBannerAll(requestId)
	responses := response.Response{Code: 200, Mssg: "success", Data: bannerResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *BannerControllerImplementation) FindBannerByDesa(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idDesa := middleware.TokenClaimsIdDesa(c)
	bannerResponse := controller.BannerServiceInterface.FindBannerByDesa(requestId, idDesa)
	responses := response.Response{Code: 200, Mssg: "success", Data: bannerResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

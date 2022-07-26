package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type PromoControllerInterface interface {
	FindPromo(c echo.Context) error
}

type PromoControllerImplementation struct {
	PromoServiceInterface service.PromoServiceInterface
}

func NewPromoController(
	promoServiceInterface service.PromoServiceInterface) PromoControllerInterface {
	return &PromoControllerImplementation{
		PromoServiceInterface: promoServiceInterface,
	}
}

func (controller *PromoControllerImplementation) FindPromo(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdDesa := middleware.TokenClaimsIdDesa(c)
	promoResponses := controller.PromoServiceInterface.FindPromo(requestId, IdDesa)
	responses := response.Response{Code: 200, Mssg: "success", Data: promoResponses, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

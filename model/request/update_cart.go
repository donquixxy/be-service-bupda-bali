package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type UpdateCartRequest struct {
	IdCart string `json:"id_cart" form:"id_cart" validate:"required"`
	Qty    int    `json:"qty" form:"qty"`
}

func ReadFromUpdateCartRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *UpdateCartRequest {
	updateCartRequest := &UpdateCartRequest{}
	if err := c.Bind(updateCartRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return updateCartRequest
}

package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type DeleteUserShippingAddressRequest struct {
	IdUserShippingAddress string `json:"id_user_shipping_address" form:"id_user_shipping_address" validate:"required"`
}

func ReadFromUserShippingAddressRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *DeleteUserShippingAddressRequest {
	deleteUserShippingAddress := &DeleteUserShippingAddressRequest{}
	if err := c.Bind(deleteUserShippingAddress); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return deleteUserShippingAddress
}

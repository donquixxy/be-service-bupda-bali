package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type CreateUserShippingAddressRequest struct {
	AlamatPengiriman string  `json:"alamat_pengiriman" form:"alamat_pengiriman" validate:"required"`
	Latitude         float64 `json:"latitude" form:"latitude" validate:"required"`
	Longitude        float64 `json:"longitude" form:"longitude" validate:"required"`
	Radius           float64 `json:"radius" form:"radius" validate:"required"`
	StatusPrimary    int     `json:"status_primary" form:"status_primary" validate:"required"`
	Catatan          string  `json:"catatan" form:"catatan" validate:"required"`
}

func ReadFromCreateUserShippingAddressRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *CreateUserShippingAddressRequest {
	createUserAddressRequest := &CreateUserShippingAddressRequest{}
	if err := c.Bind(createUserAddressRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return createUserAddressRequest
}

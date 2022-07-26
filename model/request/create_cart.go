package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type CreateCartRequest struct {
	IdProduct string `json:"id_product" form:"id_product" validate:"required"`
	Qty       int    `json:"qty" form:"qty" validate:"required"`
}

func ReadFromCreateCartRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *CreateCartRequest {
	createCartRequest := &CreateCartRequest{}
	if err := c.Bind(createCartRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return createCartRequest
}

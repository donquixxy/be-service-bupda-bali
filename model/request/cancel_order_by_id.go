package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type OrderIdRequest struct {
	IdOrder string `json:"id_order" form:"id_order" validate:"required"`
}

func ReadFromOrderIdRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *OrderIdRequest {
	orderIdRequest := &OrderIdRequest{}
	if err := c.Bind(orderIdRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return orderIdRequest
}

package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type UpdatePaymentStatusOrderRequest struct {
	TrxId       int    `json:"trx_id" form:"trx_id" validate:"required"`
	Status      string `json:"status" form:"status" validate:"required"`
	StatusCode  int    `json:"status_code" form:"status_code" validate:"required"`
	ReferenceId string `json:"reference_id" form:"reference_id" validate:"required"`
}

func ReadFromUpdatePaymentStatusOrderRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *UpdatePaymentStatusOrderRequest {
	updatePaymentStatusOrderRequest := &UpdatePaymentStatusOrderRequest{}
	if err := c.Bind(updatePaymentStatusOrderRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return updatePaymentStatusOrderRequest
}

package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type GetPaymentChannelRequest struct {
	TotalBill float64 `json:"total_bill" form:"total_bill" validate:"required"`
}

func ReadFromGetPaymentChannelRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *GetPaymentChannelRequest {
	getPaymentChannelRequest := &GetPaymentChannelRequest{}
	if err := c.Bind(getPaymentChannelRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return getPaymentChannelRequest
}

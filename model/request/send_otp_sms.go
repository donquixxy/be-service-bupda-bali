package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type SendOtpBySmsRequest struct {
	Phone string `json:"phone" form:"phone" validate:"required"`
}

func ReadFromSendOtpBySmsRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *SendOtpBySmsRequest {
	sendOtpBySmsRequest := &SendOtpBySmsRequest{}
	if err := c.Bind(sendOtpBySmsRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return sendOtpBySmsRequest
}

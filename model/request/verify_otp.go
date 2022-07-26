package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type VerifyOtpRequest struct {
	Phone   string `json:"phone" form:"phone" validate:"required"`
	OtpCode string `json:"otp_code" form:"otp_code" validate:"required"`
}

func ReadFromVerifyOtpRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *VerifyOtpRequest {
	verifyOtpRequest := &VerifyOtpRequest{}
	if err := c.Bind(verifyOtpRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return verifyOtpRequest
}

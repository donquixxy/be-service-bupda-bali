package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type InquiryPostpaidTelcoRequest struct {
	Code string `json:"code" form:"code" validate:"required"`
	Hp   string `json:"hp" form:"hp" validate:"required"`
}

func ReadFromInquiryPostpaidTelcoRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *InquiryPostpaidTelcoRequest {
	inquiryPostpaidTelcoRequest := &InquiryPostpaidTelcoRequest{}
	if err := c.Bind(inquiryPostpaidTelcoRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return inquiryPostpaidTelcoRequest
}

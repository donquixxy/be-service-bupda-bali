package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type InquiryPostpaidPdamRequest struct {
	Code string `json:"code" form:"code" validate:"required"`
	Hp   string `json:"hp" form:"hp" validate:"required"`
}

func ReadFromInquiryPostpaidPdamRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *InquiryPostpaidPdamRequest {
	inquiryPostpaidPdamRequest := &InquiryPostpaidPdamRequest{}
	if err := c.Bind(inquiryPostpaidPdamRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return inquiryPostpaidPdamRequest
}

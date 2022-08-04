package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type InquiryPostpaidPlnRequest struct {
	CustomerId string `json:"customer_id" form:"customer_id" validate:"required"`
}

func ReadFromInquiryPostpaidPlnRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *InquiryPostpaidPlnRequest {
	inquiryPostpaidPlnRequest := &InquiryPostpaidPlnRequest{}
	if err := c.Bind(inquiryPostpaidPlnRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return inquiryPostpaidPlnRequest
}

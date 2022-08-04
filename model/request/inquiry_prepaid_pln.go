package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type InquiryPrepaidPlnRequest struct {
	CustomerId string `json:"customer_id" form:"customer_id" validate:"required"`
}

func ReadFromInquiryPrepaidPlnRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *InquiryPrepaidPlnRequest {
	inquiryPrepaidPlnRequest := &InquiryPrepaidPlnRequest{}
	if err := c.Bind(inquiryPrepaidPlnRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return inquiryPrepaidPlnRequest
}

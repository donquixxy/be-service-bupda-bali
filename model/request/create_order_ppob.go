package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

// Prepaid Pulsa
type CreateOrderPrepaidRequest struct {
	ProductCode    string  `json:"product_code" form:"product_code" validate:"required"`
	CustomerId     string  `json:"customer_id" form:"customer_id" validate:"required"`
	PaymentMethod  string  `json:"payment_method" form:"payment_method" validate:"required"`
	PaymentChannel string  `json:"payment_channel" form:"payment_channel" validate:"required"`
	PaymentPoint   float64 `json:"payment_point" form:"payment_point"`
	PaymentFee     float64 `json:"payment_fee" form:"payment_fee"`
	TotalBill      float64 `json:"total_bill" form:"total_bill" validate:"required"`
}

func ReadFromCreateOrderPrepaidRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *CreateOrderPrepaidRequest {
	createOrderPrepaidRequest := &CreateOrderPrepaidRequest{}
	if err := c.Bind(createOrderPrepaidRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return createOrderPrepaidRequest
}

// Postpaid Pln
type CreateOrderPostpaidRequest struct {
	ProductCode    string  `json:"product_code" form:"product_code" validate:"required"`
	CustomerId     string  `json:"customer_id" form:"customer_id" validate:"required"`
	RefId          string  `json:"ref_id" form:"ref_id" validate:"required"`
	PaymentMethod  string  `json:"payment_method" form:"payment_method" validate:"required"`
	PaymentChannel string  `json:"payment_channel" form:"payment_channel" validate:"required"`
	PaymentPoint   float64 `json:"payment_point" form:"payment_point"`
	PaymentFee     float64 `json:"payment_fee" form:"payment_fee"`
	TotalBill      float64 `json:"total_bill" form:"total_bill" validate:"required"`
}

func ReadFromCreateOrderPostpaidRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *CreateOrderPostpaidRequest {
	createOrderPostpaidRequest := &CreateOrderPostpaidRequest{}
	if err := c.Bind(createOrderPostpaidRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return createOrderPostpaidRequest
}

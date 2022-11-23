package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type CreateOrderRequest struct {
	TotalBill        float64 `json:"total_bill" form:"total_bill" validate:"required"`
	PaymentMethod    string  `json:"payment_method" form:"payment_method" validate:"required"`
	PaymentChannel   string  `json:"payment_channel" form:"payment_channel" validate:"required"`
	PaymentPoint     float64 `json:"payment_point" form:"payment_point"`
	AlamatPengiriman string  `json:"alamat_pengiriman" form:"alamat_pengiriman" validate:"required"`
	CatatanKurir     string  `json:"catatan_kurir" form:"catatan_kurir" validate:"required"`
	ShippingCost     float64 `json:"shipping_cost" form:"shipping_cost"`
	PaymentFee       float64 `json:"payment_fee" form:"payment_fee"`
}

func ReadFromCreateOrderRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *CreateOrderRequest {
	createOrderRequest := &CreateOrderRequest{}
	if err := c.Bind(createOrderRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return createOrderRequest
}

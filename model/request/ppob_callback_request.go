package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type PpobCallbackRequest struct {
	RefId       string `json:"ref_id" form:"ref_id"`
	Status      string `json:"status" form:"status"`
	ProductCode string `json:"product_code" form:"ref_id"`
	CustomerId  string `json:"customer_id" form:"customer_id"`
	Price       string `json:"price" form:"price"`
	Message     string `json:"message" form:"message"`
	Sn          string `json:"sn" form:"sn"`
	Pin         string `json:"pin" form:"pin"`
	Balance     string `json:"balance" form:"balance"`
	TrId        string `json:"tr_id" form:"tr_id"`
	Rc          string `json:"rc" form:"rc"`
	Sign        string `json:"sign" form:"sign"`
}

func ReadFromPpobCallbackRequestRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *PpobCallbackRequest {
	ppobCallbackRequest := &PpobCallbackRequest{}
	if err := c.Bind(ppobCallbackRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return ppobCallbackRequest
}

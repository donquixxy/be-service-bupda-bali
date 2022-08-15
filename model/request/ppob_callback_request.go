package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type PpobCallbackRequest struct {
	Data PpobCallbackRequestData `json:"data"`
}

type PpobCallbackRequestData struct {
	RefId       string `json:"ref_id"`
	Status      string `json:"status"`
	ProductCode string `json:"product_code"`
	CustomerId  string `json:"customer_id"`
	Price       string `json:"price"`
	Message     string `json:"message"`
	Sn          string `json:"sn"`
	Pin         string `json:"pin"`
	Balance     string `json:"balance"`
	TrId        string `json:"tr_id"`
	Rc          string `json:"rc"`
	Sign        string `json:"sign"`
}

func ReadFromPpobCallbackRequestRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *PpobCallbackRequest {
	ppobCallbackRequest := &PpobCallbackRequest{}
	if err := c.Bind(ppobCallbackRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return ppobCallbackRequest
}

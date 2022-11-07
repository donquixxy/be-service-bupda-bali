package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type CreatePaylaterRequest struct {
	Amount float64 `json:"amount"`
}

func ReadFromCreatePaylaterRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *CreatePaylaterRequest {
	createPaylaterRequest := &CreatePaylaterRequest{}
	if err := c.Bind(createPaylaterRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return createPaylaterRequest
}

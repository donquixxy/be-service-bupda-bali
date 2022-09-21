package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type FindBigisResponsesRequest struct {
	Nik string `json:"nik" form:"nik" validate:"required"`
}

func ReadFromFindBigisResponsesRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *FindBigisResponsesRequest {
	findBigisResponsesRequest := &FindBigisResponsesRequest{}
	if err := c.Bind(findBigisResponsesRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return findBigisResponsesRequest
}

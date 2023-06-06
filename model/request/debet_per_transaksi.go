package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type DebetPerTransaksiRequest struct {
	LoanId []string `json:"loan_id" form:"loan_id" validate:"required"`
}

func ReadFromDebetPerTransaksiRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *DebetPerTransaksiRequest {
	debetPerTransaksiRequest := &DebetPerTransaksiRequest{}
	if err := c.Bind(debetPerTransaksiRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return debetPerTransaksiRequest
}

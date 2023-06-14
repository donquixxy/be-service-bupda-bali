package request

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type DebetPerTransaksiRequest struct {
	LoanId []string `json:"loan_id" form:"loan_id[]" validate:"required"`
}

func ReadFromDebetPerTransaksiRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *DebetPerTransaksiRequest {
	debetPerTransaksiRequest := &DebetPerTransaksiRequest{}

	count := 0
	for {
		index := fmt.Sprintf("loan_id[%d]", count)
		if value := c.FormValue(index); value != "" {
			debetPerTransaksiRequest.LoanId = append(debetPerTransaksiRequest.LoanId, value)
			count++
		} else {
			break
		}
	}

	log.Println("...any", debetPerTransaksiRequest)

	return debetPerTransaksiRequest
}

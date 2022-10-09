package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type CreateMerchantApproveListRequest struct {
	NamaMerchant string `json:"nama_merchant" form:"nama_merchant" validate:"required"`
}

func ReadFromCreateMerchantApproveListRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *CreateMerchantApproveListRequest {
	createMerchantApproveListRequest := &CreateMerchantApproveListRequest{}
	if err := c.Bind(createMerchantApproveListRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return createMerchantApproveListRequest
}

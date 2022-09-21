package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type UpdateUserPhoneRequest struct {
	Phone     string `json:"phone" form:"phone" validate:"required"`
	FormToken string `json:"form_token" form:"form_token" validate:"required"`
}

func ReadFromUpdateUserPhoneRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *UpdateUserPhoneRequest {
	updateUserPhoneRequest := &UpdateUserPhoneRequest{}
	if err := c.Bind(updateUserPhoneRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return updateUserPhoneRequest
}

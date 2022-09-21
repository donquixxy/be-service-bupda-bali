package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type UpdateUserProfileRequest struct {
	NamaLengkap string `json:"nama_lengkap" form:"nama_lengkap" validate:"required"`
	Email       string `json:"email" form:"email" validate:"required"`
}

func ReadFromUpdateUserProfileRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *UpdateUserProfileRequest {
	updateUserProfileRequest := &UpdateUserProfileRequest{}
	if err := c.Bind(updateUserProfileRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return updateUserProfileRequest
}

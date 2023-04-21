package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type UpdateUserPasswordRequest struct {
	PasswordLama string `json:"password_lama" form:"password_lama" validate:"required"`
	PasswordBaru string `json:"password_baru" form:"password_baru" validate:"required"`
}

func ReadFromUpdateUserPasswordRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *UpdateUserPasswordRequest {
	updateUserPasswordRequest := &UpdateUserPasswordRequest{}
	if err := c.Bind(updateUserPasswordRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return updateUserPasswordRequest
}

type UpdateUserForgotPasswordRequest struct {
	Phone        string `json:"phone" form:"phone" validate:"required"`
	PasswordBaru string `json:"password_baru" form:"password_baru" validate:"required"`
	FormToken    string `json:"form_token" form:"form_token" validate:"required"`
}

func ReadFromUpdateUserForgotPasswordRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *UpdateUserForgotPasswordRequest {
	updateUserForgotPasswordRequest := &UpdateUserForgotPasswordRequest{}
	if err := c.Bind(updateUserForgotPasswordRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return updateUserForgotPasswordRequest
}

type UpdateUserPasswordInveliRequest struct {
	Phone       string `json:"phone" form:"phone" validate:"required"`
	NewPassword string `json:"new_password" form:"new_password" validate:"required"`
}

func ReadFromUpdateUserPasswordInveliRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *UpdateUserPasswordInveliRequest {
	updateUserPasswordInveliRequest := &UpdateUserPasswordInveliRequest{}
	if err := c.Bind(updateUserPasswordInveliRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return updateUserPasswordInveliRequest
}

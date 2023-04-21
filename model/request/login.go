package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type LoginRequest struct {
	Phone    string `json:"phone" form:"phone" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

func ReadFromLoginRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *LoginRequest {
	loginRequest := &LoginRequest{}
	if err := c.Bind(loginRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return loginRequest
}

type UbahPasswordInveliRequest struct {
	Phone              string `json:"phone" form:"phone" validate:"required"`
	PasswordFromInveli string `json:"password_from_inveli" form:"password_from_inveli" validate:"required"`
	NewPassword        string `json:"new_password" form:"new_password" validate:"required"`
}

func ReadFromLoginInveliRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *UbahPasswordInveliRequest {
	ubahPasswordInveliRequest := &UbahPasswordInveliRequest{}
	if err := c.Bind(ubahPasswordInveliRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return ubahPasswordInveliRequest
}

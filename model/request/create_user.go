package request

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

type CreateUserRequest struct {
	NoIdentitas string `json:"no_identitas" form:"no_identitas" validate:"required"`
	NamaLengkap string `json:"nama_lengkap" form:"nama_lengkap" validate:"required"`
	Phone       string `json:"phone" form:"phone" validate:"required"`
	Email       string `json:"email" form:"email"`
	IdDesa      string `json:"id_desa" form:"id_desa" validate:"required"`
	// Password    string `json:"password" form:"password" validate:"required"`
	FormToken string `json:"form_token" form:"form_token" validate:"required"`
}

func ReadFromCreateUserRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *CreateUserRequest {
	createUserRequest := &CreateUserRequest{}
	if err := c.Bind(createUserRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return createUserRequest
}

type CreateUserSurveyedRequest struct {
	NoIdentitas string `json:"no_identitas" form:"no_identitas" validate:"required"`
	NamaLengkap string `json:"nama_lengkap" form:"nama_lengkap" validate:"required"`
	Phone       string `json:"phone" form:"phone" validate:"required"`
	Email       string `json:"email" form:"email"`
	IdDesa      string `json:"id_desa" form:"id_desa" validate:"required"`
	// Password    string `json:"password" form:"password" validate:"required"`
	Alamat string `json:"alamat" form:"alamat" validate:"required"`
}

func ReadFromCreateUserSurveyedRequestBody(c echo.Context, requestId string, logger *logrus.Logger) *CreateUserSurveyedRequest {
	createUserRequest := &CreateUserSurveyedRequest{}
	if err := c.Bind(createUserRequest); err != nil {
		exceptions.PanicIfError(err, requestId, logger)
	}
	return createUserRequest
}

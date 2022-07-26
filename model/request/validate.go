package request

import (
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
)

func ValidateRequest(validate *validator.Validate, request interface{}, requestId string, logger *logrus.Logger) {
	var errorStrings []string
	err := validate.Struct(request)
	var errorString string
	if err != nil {
		for _, errorValidation := range err.(validator.ValidationErrors) {
			errorString = errorValidation.Field() + " is " + errorValidation.Tag()
			errorStrings = append(errorStrings, errorString)
		}
		exceptions.PanicIfBadRequest(err, requestId, errorStrings, logger)
	}
}

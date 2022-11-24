package exceptions

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ErrorStruct struct {
	Code  int      `json:"code"`
	Mssg  string   `json:"message"`
	Error []string `json:"error"`
}

func PanicIfError(err error, requestId string, logger *logrus.Logger) {
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.WithFields(logrus.Fields{"request_id": requestId}).Error(err)
		panic(err)
	}
}

func PanicIfErrorWithRollback(err error, requestId string, errorString []string, logger *logrus.Logger, DB *gorm.DB) {
	if err != nil && err != gorm.ErrRecordNotFound {
		rollback := DB.Rollback()
		if rollback.Error != nil {
			PanicIfError(rollback.Error, requestId, logger)
		}
		out, errr := json.Marshal(ErrorStruct{Code: 404, Error: errorString})
		PanicIfError(errr, requestId, logger)
		logger.WithFields(logrus.Fields{"request_id": requestId}).Error(err)
		panic(string(out))
	}
}

func PanicIfErrorWithRollbackRegister(err error, requestId string, errorString []string, logger *logrus.Logger, DB *gorm.DB) {
	if err != nil && err != gorm.ErrRecordNotFound {
		rollback := DB.Rollback()
		if rollback.Error != nil {
			PanicIfError(rollback.Error, requestId, logger)
		}
		out, errr := json.Marshal(ErrorStruct{Code: 408, Error: errorString})
		PanicIfError(errr, requestId, logger)
		logger.WithFields(logrus.Fields{"request_id": requestId}).Error(err)
		panic(string(out))
	}
}

func PanicIfBadRequest(err error, requestId string, errorString []string, logger *logrus.Logger) {
	if err != nil {
		out, errr := json.Marshal(ErrorStruct{Code: 400, Mssg: "Bad Request", Error: errorString})
		PanicIfError(errr, requestId, logger)
		logger.WithFields(logrus.Fields{"request_id": requestId}).Error(err)
		panic(string(out))
	}
}

func PanicIfRecordNotFound(err error, requestId string, errorString []string, logger *logrus.Logger) {
	if err != nil || err == gorm.ErrRecordNotFound {
		out, errr := json.Marshal(ErrorStruct{Code: 404, Mssg: "Not Found", Error: errorString})
		PanicIfError(errr, requestId, logger)
		logger.WithFields(logrus.Fields{"request_id": requestId}).Error(err)
		panic(string(out))
	}
}

func PanicIfRecordNotFoundWithRollback(err error, requestId string, errorString []string, logger *logrus.Logger, DB *gorm.DB) {
	if err != nil && err == gorm.ErrRecordNotFound {
		rollback := DB.Rollback()
		if rollback.Error != nil {
			PanicIfError(rollback.Error, requestId, logger)
		}
		out, errr := json.Marshal(ErrorStruct{Code: 404, Error: errorString})
		PanicIfError(errr, requestId, logger)
		logger.WithFields(logrus.Fields{"request_id": requestId}).Error(err)
		panic(string(out))
	}
}

func PanicIfUnauthorized(err error, requestId string, errorString []string, logger *logrus.Logger) {
	if err != nil {
		out, errr := json.Marshal(ErrorStruct{Code: 401, Mssg: "Unauthorized", Error: errorString})
		PanicIfError(errr, requestId, logger)
		logger.WithFields(logrus.Fields{"request_id": requestId}).Error(err)
		panic(string(out))
	}
}

func PanicIfRecordAlreadyExists(err error, requestId string, errorString []string, logger *logrus.Logger) {
	if err != nil {
		out, errr := json.Marshal(ErrorStruct{Code: 409, Mssg: "Data Already Exist", Error: errorString})
		PanicIfError(errr, requestId, logger)
		logger.WithFields(logrus.Fields{"request_id": requestId}).Error(err)
		panic(string(out))
	}
}

func PanicIfRecordAlreadyExistsWIthRollback(err error, requestId string, errorString []string, logger *logrus.Logger, DB *gorm.DB) {
	if err != nil {
		rollback := DB.Rollback()
		if rollback.Error != nil {
			PanicIfError(rollback.Error, requestId, logger)
		}
		out, errr := json.Marshal(ErrorStruct{Code: 409, Error: errorString})
		PanicIfError(errr, requestId, logger)
		logger.WithFields(logrus.Fields{"request_id": requestId}).Error(err)
		panic(string(out))
	}
}

func PanicPPOBHandler(err error, requestId string, message string, errorString []string, logger *logrus.Logger) {
	if err != nil {
		out, errr := json.Marshal(ErrorStruct{Code: 409, Mssg: message, Error: errorString})
		PanicIfError(errr, requestId, logger)
		logger.WithFields(logrus.Fields{"request_id": requestId}).Error(err)
		panic(string(out))
	}
}

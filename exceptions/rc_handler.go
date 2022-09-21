package exceptions

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
)

type RcStruct struct {
	Code int    `json:"code"`
	Mssg string `json:"message"`
	Rc   int    `json:"rc"`
}

func RcHandler(err error, e echo.Context) {
	errS := ErrorStruct{}
	json.Unmarshal([]byte(err.Error()), &errS)
	if errS.Code != 0 {
		response := response.Response{Code: errS.Code, Mssg: errS.Mssg, Data: []string{}, Error: errS.Error}
		e.JSON(errS.Code, response)
	} else {
		response := response.Response{Data: []string{}, Error: []string{"Internal Server Error"}}
		e.JSON(http.StatusInternalServerError, response)
	}
}

func RcResponse(err error, requestId string, errorString []string, logger *logrus.Logger) {
	if err != nil {
		out, errr := json.Marshal(ErrorStruct{Code: 400, Mssg: "Bad Request", Error: errorString})
		PanicIfError(errr, requestId, logger)
		logger.WithFields(logrus.Fields{"request_id": requestId}).Error(err)
		panic(string(out))
	}
}

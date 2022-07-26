package exceptions

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
)

func ErrorHandler(err error, e echo.Context) {
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

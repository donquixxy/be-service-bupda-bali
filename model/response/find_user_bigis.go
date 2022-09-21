package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/bigis"

type FindUserFromBigisResponse struct {
	Nik    string `json:"nik"`
	Name   string `json:"name"`
	Alamat string `json:"alamat"`
	Phone  string `json:"phone"`
}

func ToFindUserFromBigisResponse(responses *bigis.Response) (bigisResponse FindUserFromBigisResponse) {
	bigisResponse.Nik = responses.DataResponse.Nik
	bigisResponse.Name = responses.DataResponse.Name
	bigisResponse.Alamat = responses.DataResponse.Alamat
	bigisResponse.Phone = responses.DataResponse.Phone
	return bigisResponse
}

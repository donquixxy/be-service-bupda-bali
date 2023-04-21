package response

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/bigis"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type FindUserFromBigisResponse struct {
	Nik       string `json:"nik"`
	Name      string `json:"name"`
	Alamat    string `json:"alamat"`
	Kelurahan string `json:"kelurahan"`
	Kecamatan string `json:"kecamatan"`
	Desa      string `json:"desa"`
	IdDesa    string `json:"id_desa"`
}

func ToFindUserFromBigisResponse(responses *bigis.Response, desa *entity.Desa) (bigisResponse FindUserFromBigisResponse) {
	bigisResponse.Nik = responses.DataResponse.Nik
	bigisResponse.Name = responses.DataResponse.Name
	bigisResponse.Alamat = responses.DataResponse.Alamat
	bigisResponse.Kelurahan = responses.DataResponse.Kelurahan
	bigisResponse.Kecamatan = responses.DataResponse.Kecamatan
	bigisResponse.Desa = desa.NamaDesa
	bigisResponse.IdDesa = desa.Id
	return bigisResponse
}

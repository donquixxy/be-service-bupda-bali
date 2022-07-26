package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindKecamatanResponse struct {
	IdKeca   int    `json:"idkeca"`
	IdKabu   int    `json:"kdkabu"`
	IdProp   int    `json:"idprop"`
	KdKeca   string `json:"kdkeca"`
	NamaKeca string `json:"nama_keca"`
}

func ToFindKecamatanResponse(kecamatans []entity.Kecamatan) (kecamatanResponses []FindKecamatanResponse) {
	for _, kecamatan := range kecamatans {
		kecamatanResponse := FindKecamatanResponse{}
		kecamatanResponse.IdKeca = kecamatan.IdKeca
		kecamatanResponse.IdKabu = kecamatan.IdKabu
		kecamatanResponse.IdProp = kecamatan.IdProp
		kecamatanResponse.KdKeca = kecamatan.KdKeca
		kecamatanResponse.NamaKeca = kecamatan.NamaKeca
		kecamatanResponses = append(kecamatanResponses, kecamatanResponse)
	}
	return kecamatanResponses
}

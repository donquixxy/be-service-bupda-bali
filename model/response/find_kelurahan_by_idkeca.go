package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindKelurahanByIdKecaResponse struct {
	IdKelu   int    `json:"idkelu"`
	IdKeca   int    `json:"kdkeca"`
	IdKabu   int    `json:"idkabu"`
	IdProp   int    `json:"idprop"`
	KdKelu   string `json:"kdkelu"`
	NamaKelu string `json:"nama_kelu"`
}

func ToFindKelurahanByIdKecaResponse(kelurahans []entity.Kelurahan) (kelurahanResponses []FindKelurahanByIdKecaResponse) {
	for _, kelurahan := range kelurahans {
		var kelurahanResponse FindKelurahanByIdKecaResponse
		kelurahanResponse.IdKelu = kelurahan.IdKelu
		kelurahanResponse.IdKeca = kelurahan.IdKeca
		kelurahanResponse.IdKabu = kelurahan.IdKabu
		kelurahanResponse.IdProp = kelurahan.IdProp
		kelurahanResponse.KdKelu = kelurahan.KdKelu
		kelurahanResponse.NamaKelu = kelurahan.NamaKelu
		kelurahanResponses = append(kelurahanResponses, kelurahanResponse)
	}
	return kelurahanResponses
}

package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindDesaByIdKeluResponse struct {
	Id          string `json:"id"`
	NamaDesa    string `json:"nama_desa"`
	NamaBendesa string `json:"nama_bendesa"`
}

func ToFindDesaByIdKeluResponse(desas []entity.Desa) (desasResponses []FindDesaByIdKeluResponse) {
	for _, desa := range desas {
		var desaResponse FindDesaByIdKeluResponse
		desaResponse.Id = desa.Id
		desaResponse.NamaDesa = desa.NamaDesa
		desaResponse.NamaBendesa = desa.NamaBendesa
		desasResponses = append(desasResponses, desaResponse)
	}
	return desasResponses
}

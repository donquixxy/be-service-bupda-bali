package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindSettingShippingCostResponse struct {
	Id          string  `json:"id"`
	IdDesa      string  `json:"id_desa"`
	SettingName string  `json:"setting_name"`
	Value       float64 `json:"shipping_cost"`
}

func ToFindSettingShippingCostResponse(setting *entity.Setting) (settingResponse FindSettingShippingCostResponse) {
	settingResponse.Id = setting.Id
	settingResponse.IdDesa = setting.IdDesa
	settingResponse.SettingName = setting.SettingName
	settingResponse.Value = setting.Value
	return settingResponse
}

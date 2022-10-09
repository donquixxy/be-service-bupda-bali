package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindVersionResponse struct {
	SettingTitle string `json:"os"`
	SettingName  string `json:"version"`
}

func ToFindVersionResponse(setting *entity.Setting) (settingResponse FindVersionResponse) {
	settingResponse.SettingTitle = setting.SettingTitle
	settingResponse.SettingName = setting.SettingName
	return settingResponse
}

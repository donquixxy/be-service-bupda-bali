package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindVersionResponse struct {
	OSName  string `json:"os"`
	Current string `json:"current"`
	New     string `json:"new"`
}

func ToFindNewVersionResponse(appVersion []entity.AppVersion, os int) (settingResponse FindVersionResponse) {
	if os == 1 {
		settingResponse.OSName = "Android"
	} else {

	}
	settingResponse.OSName = appVersion[0].OS
	settingResponse.Current = appVersion[0].Version
	settingResponse.New = appVersion[1].Version
	return settingResponse
}

func ToFindNewVersion2Response(setting []entity.Setting, os int) (settingResponse FindVersionResponse) {
	if os == 1 {
		settingResponse.OSName = "Android"
	} else {
		settingResponse.OSName = "iOS"
	}
	settingResponse.Current = setting[0].SettingName
	settingResponse.New = setting[1].SettingName
	return settingResponse
}

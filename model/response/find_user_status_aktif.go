package response

type FindUserStatusAktifResponse struct {
	StatusAktif bool `json:"status_aktif"`
}

func ToFindUserStatusAktifResponse(statusAktif int) (userResponse FindUserStatusAktifResponse) {
	if statusAktif == 2 {
		userResponse.StatusAktif = true
	} else {
		userResponse.StatusAktif = false
	}

	return userResponse
}

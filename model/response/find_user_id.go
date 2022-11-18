package response

import (
	"log"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type FindUserIdResponse struct {
	Id                string `json:"id"`
	NoIdentitas       string `json:"no_identitas"`
	NamaLengkap       string `json:"nama_lengkap"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	StatusSurvey      int    `json:"status_survey"`
	AccountType       int    `json:"account_type"`
	MerchantCode      string `json:"merchant_code"`
	IdLimitPayLater   string `json:"id_limit_pay_later"`
	NamaDesa          string `json:"nama_desa"`
	NamaBendesa       string `json:"nama_bendesa"`
	StatusAktifInveli int    `json:"status_aktif_inveli"`
	StatusPaylater    int    `json:"status_paylater"`
}

func ToFindUserIdResponse(userProfile *entity.UserProfile, statusAktifUser bool) (userResponse FindUserIdResponse) {
	log.Println("status aktif user", userProfile.User.StatusPaylater)
	userResponse.Id = userProfile.User.Id
	userResponse.NoIdentitas = userProfile.NoIdentitas
	userResponse.NamaLengkap = userProfile.NamaLengkap
	userResponse.Email = userProfile.Email
	userResponse.Phone = userProfile.User.Phone
	userResponse.StatusSurvey = userProfile.User.StatusSurvey
	userResponse.AccountType = userProfile.User.AccountType
	userResponse.MerchantCode = userProfile.User.MerchantCode
	userResponse.IdLimitPayLater = userProfile.User.IdLimitPayLater
	userResponse.NamaDesa = userProfile.User.Desa.NamaDesa
	userResponse.NamaBendesa = userProfile.User.Desa.NamaBendesa
	userResponse.StatusAktifInveli = userProfile.User.StatusPaylater
	userResponse.StatusPaylater = userProfile.User.IsPaylater
	return userResponse
}

package response

import (
	"time"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type FindMerchantStatusApproveByUserResponse struct {
	NamaLengkap   string    `json:"nama_lengkap"`
	MerchantName  string    `json:"nama_merchant"`
	ApproveStatus int       `json:"status_approve"`
	CreatedAt     time.Time `json:"tgl_pengajuan"`
}

func ToFindMerchantStatusApproveByUserResponse(merchantApproveList *entity.MerchantApproveList) (merchantResponse FindMerchantStatusApproveByUserResponse) {
	merchantResponse.NamaLengkap = merchantApproveList.NamaLengkap
	merchantResponse.MerchantName = merchantApproveList.MerchantName
	merchantResponse.ApproveStatus = merchantApproveList.ApproveStatus
	merchantResponse.CreatedAt = merchantApproveList.CreatedAt
	return merchantResponse
}

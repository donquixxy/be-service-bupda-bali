package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindPointByUserResponse struct {
	JmlPoint    float64 `json:"jml_point"`
	StatusPoint int     `json:"status_point"`
}

func ToFindPointByUserResponse(point *entity.Point) (pointResponse FindPointByUserResponse) {
	pointResponse.JmlPoint = point.JmlPoint
	pointResponse.StatusPoint = point.StatusPoint
	return pointResponse
}

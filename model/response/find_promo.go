package response

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type FindPromoResponse struct {
	Id          string `json:"id"`
	PromoTitle  string `json:"promo_title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

func ToFindPromoResponse(promos []entity.Promo) (promoResponses []FindPromoResponse) {
	for _, promo := range promos {
		var promoResponse FindPromoResponse
		promoResponse.Id = promo.Id
		promoResponse.PromoTitle = promo.PromoTitle
		promoResponse.Description = promo.Description
		promoResponse.Image = promo.Image
		promoResponse.StartDate = promo.StartDate.Format("2006-01-02")
		promoResponse.EndDate = promo.EndDate.Format("2006-01-02")
		promoResponses = append(promoResponses, promoResponse)
	}
	return promoResponses
}

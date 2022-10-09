package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindBannerAllResponse struct {
	Id              string `json:"id"`
	IdDesa          string `json:"id_desa"`
	BannerTitle     string `json:"banner_title"`
	BannerImg       string `json:"banner_img"`
	BannerUrl       string `json:"banner_url"`
	BannerReference string `json:"banner_reference"`
}

func ToFindBannerAllResponse(banners []entity.Banner) (bannersResponses []FindBannerAllResponse) {
	for _, banner := range banners {
		var bannerResponse FindBannerAllResponse
		bannerResponse.Id = banner.Id
		bannerResponse.IdDesa = banner.IdDesa
		bannerResponse.BannerTitle = banner.BannerTitle
		bannerResponse.BannerImg = banner.BannerImg
		bannerResponse.BannerUrl = banner.BannerUrl
		bannerResponse.BannerReference = banner.BannerReference
		bannersResponses = append(bannersResponses, bannerResponse)
	}
	return bannersResponses
}

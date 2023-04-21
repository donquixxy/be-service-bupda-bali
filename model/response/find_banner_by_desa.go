package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindBannerByIdDesaResponse struct {
	Id              string `json:"id"`
	IdDesa          string `json:"id_desa"`
	BannerTitle     string `json:"banner_title"`
	BannerImg       string `json:"banner_img"`
	BannerUrl       string `json:"banner_url"`
	BannerReference string `json:"banner_reference"`
}

func ToFindBannerByIdDesaResponse(banners []entity.Banner) (bannersResponses []FindBannerByIdDesaResponse) {
	for _, banner := range banners {
		var bannerResponse FindBannerByIdDesaResponse
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

package response

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type FindProductsDesaResponse struct {
	Id              string  `json:"id"`
	IdBrand         int     `json:"id_brand"`
	IdCategory      int     `json:"id_category"`
	IdSubCategory   int     `json:"id_sub_category"`
	IdUnit          int     `json:"id_unit"`
	NoSku           string  `json:"no_sku"`
	ProductName     string  `json:"product_name"`
	Price           float64 `json:"price"`
	PromoPercentage float64 `json:"promo_percentage"`
	FlagPromo       int     `json:"flag_promo"`
	Description     string  `json:"description"`
	PictureUrl      string  `json:"picture_url"`
	Thumbnail       string  `json:"thumbnail"`
	StockOpname     int     `json:"stock_opname"`
	PriceInfo       string  `json:"price_info"`
	AccountType     string  `json:"account_type"`
}

func ToFindProductsDesaResponse(productsDesas []entity.ProductsDesa, AccountType int) (productsDesaResponses []FindProductsDesaResponse) {
	for _, productDesa := range productsDesas {
		productsDesaResponse := FindProductsDesaResponse{}
		productsDesaResponse.Id = productDesa.Id
		productsDesaResponse.IdBrand = productDesa.ProductsMaster.IdBrand
		productsDesaResponse.IdCategory = productDesa.ProductsMaster.IdCategory
		productsDesaResponse.IdSubCategory = productDesa.ProductsMaster.IdSubCategory
		productsDesaResponse.IdUnit = productDesa.ProductsMaster.IdUnit
		productsDesaResponse.NoSku = productDesa.ProductsMaster.NoSku
		productsDesaResponse.ProductName = productDesa.ProductsMaster.ProductName
		productsDesaResponse.FlagPromo = productDesa.IsPromo
		if AccountType == 1 {
			if productDesa.IsPromo == 1 {
				productsDesaResponse.AccountType = "User Biasa"
				productsDesaResponse.PriceInfo = "Harga Normal"
				productsDesaResponse.Price = productDesa.PricePromo
				productsDesaResponse.PromoPercentage = productDesa.PercentagePromo
			} else {
				productsDesaResponse.AccountType = "User Biasa"
				productsDesaResponse.PriceInfo = "Harga Normal"
				productsDesaResponse.Price = productDesa.ProductsMaster.Price
			}
		} else if AccountType == 2 {
			productsDesaResponse.AccountType = "User Merchant"
			productsDesaResponse.PriceInfo = "Harga Grosir"
			productsDesaResponse.Price = productDesa.ProductsMaster.PriceGrosir
		}
		productsDesaResponse.Description = productDesa.ProductsMaster.Description
		productsDesaResponse.PictureUrl = productDesa.ProductsMaster.PictureUrl
		productsDesaResponse.Thumbnail = productDesa.ProductsMaster.Thumbnail
		productsDesaResponse.Thumbnail = productDesa.ProductsMaster.Thumbnail
		productsDesaResponse.StockOpname = productDesa.StockOpname
		productsDesaResponses = append(productsDesaResponses, productsDesaResponse)
	}
	return productsDesaResponses
}

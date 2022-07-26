package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindProductDesaByIdResponse struct {
	Id              string  `json:"id"`
	IdBrand         int     `json:"id_brand"`
	IdCategory      int     `json:"id_category"`
	IdSubCategory   int     `json:"id_sub_category"`
	IdUnit          int     `json:"id_unit"`
	NoSku           string  `json:"no_sku"`
	ProductName     string  `json:"product_name"`
	Price           float64 `json:"price"`
	PromoPercentage float64 `json:"promo_percentage"`
	IsPromo         int     `json:"flag_promo"`
	Description     string  `json:"description"`
	PictureUrl      string  `json:"picture_url"`
	Thumbnail       string  `json:"thumbnail"`
	StockOpname     int     `json:"stock_opname"`
	PriceInfo       string  `json:"price_info"`
	AccountType     string  `json:"account_type"`
}

func ToFindProductDesaByIdResponse(productDesa *entity.ProductsDesa, AccountType int) (productDesaResponse FindProductDesaByIdResponse) {
	productDesaResponse.Id = productDesa.Id
	productDesaResponse.IdBrand = productDesa.ProductsMaster.IdBrand
	productDesaResponse.IdCategory = productDesa.ProductsMaster.IdCategory
	productDesaResponse.IdSubCategory = productDesa.ProductsMaster.IdSubCategory
	productDesaResponse.IdUnit = productDesa.ProductsMaster.IdUnit
	productDesaResponse.NoSku = productDesa.ProductsMaster.NoSku
	productDesaResponse.ProductName = productDesa.ProductsMaster.ProductName
	productDesaResponse.IsPromo = productDesa.IsPromo
	if AccountType == 1 {
		if productDesa.IsPromo == 1 {
			productDesaResponse.AccountType = "User Biasa"
			productDesaResponse.PriceInfo = "Harga Normal"
			productDesaResponse.Price = productDesa.PricePromo
			productDesaResponse.PromoPercentage = productDesa.PercentagePromo
		} else {
			productDesaResponse.AccountType = "User Biasa"
			productDesaResponse.PriceInfo = "Harga Normal"
			productDesaResponse.Price = productDesa.ProductsMaster.Price
		}
	} else if AccountType == 2 {
		productDesaResponse.AccountType = "User Merchant"
		productDesaResponse.PriceInfo = "Harga Grosir"
		productDesaResponse.Price = productDesa.ProductsMaster.PriceGrosir
	}
	productDesaResponse.Description = productDesa.ProductsMaster.Description
	productDesaResponse.PictureUrl = productDesa.ProductsMaster.PictureUrl
	productDesaResponse.Thumbnail = productDesa.ProductsMaster.Thumbnail
	productDesaResponse.StockOpname = productDesa.StockOpname
	return productDesaResponse
}

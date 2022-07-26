package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindCartByUserResponse struct {
	Id              string  `json:"id"`
	IdProduct       string  `json:"id_product"`
	ProductName     string  `json:"product_name"`
	PictureUrl      string  `json:"picture_url"`
	Thumbnail       string  `json:"thumbnail"`
	Price           float64 `json:"price"`
	PromoPercentage float64 `json:"promo_percentage"`
	IsPromo         int     `json:"is_promo"`
	Qty             int     `json:"qty"`
	Stock           int     `json:"stock"`
	Description     string  `json:"description"`
	PriceInfo       string  `json:"price_info"`
	AccountType     string  `json:"account_type"`
}

func ToFindCartByUserResponse(carts []entity.Cart, AccountType int) (cartResponses []FindCartByUserResponse) {
	for _, cart := range carts {
		cartResponse := FindCartByUserResponse{}
		cartResponse.Id = cart.Id
		cartResponse.IdProduct = cart.IdProductDesa
		cartResponse.ProductName = cart.ProductsDesa.ProductsMaster.ProductName
		cartResponse.PictureUrl = cart.ProductsDesa.ProductsMaster.PictureUrl
		cartResponse.Thumbnail = cart.ProductsDesa.ProductsMaster.Thumbnail
		cartResponse.Qty = cart.Qty
		cartResponse.Stock = cart.ProductsDesa.StockOpname
		cartResponse.Description = cart.ProductsDesa.ProductsMaster.Description
		cartResponse.IsPromo = cart.ProductsDesa.IsPromo
		if AccountType == 1 {
			if cart.ProductsDesa.IsPromo == 1 {
				cartResponse.AccountType = "User Biasa"
				cartResponse.PriceInfo = "Harga Normal"
				cartResponse.Price = cart.ProductsDesa.PricePromo
				cartResponse.PromoPercentage = cart.ProductsDesa.PercentagePromo
			} else {
				cartResponse.AccountType = "User Biasa"
				cartResponse.PriceInfo = "Harga Normal"
				cartResponse.Price = cart.ProductsDesa.ProductsMaster.Price
			}
		} else if AccountType == 2 {
			cartResponse.AccountType = "User Merchant"
			cartResponse.PriceInfo = "Harga Grosir"
			cartResponse.Price = cart.ProductsDesa.ProductsMaster.PriceGrosir
		}
		cartResponses = append(cartResponses, cartResponse)
	}
	return cartResponses
}

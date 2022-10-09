package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindCartByIdUserResponse struct {
	SubTotal     float64    `json:"sub_total"`
	ShippingCost float64    `json:"shipping_cost"`
	TotalBill    float64    `json:"total_bill"`
	CartItems    []CartItem `json:"cart_items"`
}

type CartItem struct {
	Id              string  `json:"id"`
	IdProduct       string  `json:"id_product"`
	ProductName     string  `json:"product_name"`
	PictureUrl      string  `json:"picture_url"`
	Thumbnail       string  `json:"thumbnail"`
	Price           float64 `json:"price_normal"`
	PricePromo      float64 `json:"price_promo"`
	PromoPercentage float64 `json:"promo_percentage"`
	IsPromo         int     `json:"is_promo"`
	Qty             int     `json:"qty"`
	Stock           int     `json:"stock"`
	Description     string  `json:"description"`
	PriceInfo       string  `json:"price_info"`
	AccountType     string  `json:"account_type"`
}

func ToFindCartByUserResponse(carts []entity.Cart, ShippingCost float64, AccountType int) (cartResponse FindCartByIdUserResponse) {
	var cartItems []CartItem
	var totalPricePerItem float64
	var subTotal float64
	for _, cart := range carts {
		var cartItem CartItem
		cartItem.Id = cart.Id
		cartItem.IdProduct = cart.IdProductDesa
		cartItem.ProductName = cart.ProductsDesa.ProductsMaster.ProductName
		cartItem.PictureUrl = cart.ProductsDesa.ProductsMaster.PictureUrl
		cartItem.Thumbnail = cart.ProductsDesa.ProductsMaster.Thumbnail
		cartItem.Qty = cart.Qty
		cartItem.Stock = cart.ProductsDesa.StockOpname
		cartItem.Description = cart.ProductsDesa.ProductsMaster.Description
		cartItem.IsPromo = cart.ProductsDesa.IsPromo
		if AccountType == 1 {
			if cart.ProductsDesa.IsPromo == 1 {
				cartItem.PricePromo = cart.ProductsDesa.PricePromo
				cartItem.Price = cart.ProductsDesa.Price
				totalPricePerItem = cart.ProductsDesa.PricePromo * float64(cart.Qty)
				cartItem.PromoPercentage = cart.ProductsDesa.PercentagePromo
			} else {
				cartItem.Price = cart.ProductsDesa.Price
				totalPricePerItem = cart.ProductsDesa.Price * float64(cart.Qty)
				cartItem.PricePromo = 0
				cartItem.PromoPercentage = 0
			}
			cartItem.AccountType = "User Biasa"
			cartItem.PriceInfo = "Krama Harga Normal"
		} else if AccountType == 2 {
			cartItem.AccountType = "User Merchant"
			cartItem.PriceInfo = "Krama Harga Grosir"
			cartItem.Price = cart.ProductsDesa.PriceGrosir
			totalPricePerItem = cart.ProductsDesa.PriceGrosir * float64(cart.Qty)
		}
		subTotal = subTotal + totalPricePerItem
		cartItems = append(cartItems, cartItem)
	}

	cartResponse.CartItems = cartItems
	cartResponse.SubTotal = subTotal
	cartResponse.ShippingCost = ShippingCost
	cartResponse.TotalBill = subTotal + ShippingCost

	return cartResponse
}

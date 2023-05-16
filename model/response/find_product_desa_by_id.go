package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindProductDesaByIdResponse struct {
	Id               string             `json:"id"`
	IdBrand          int                `json:"id_brand"`
	IdCategory       int                `json:"id_category"`
	IdSubCategory    int                `json:"id_sub_category"`
	IdType           int                `json:"id_type"`
	IdUnit           int                `json:"id_unit"`
	NoSku            string             `json:"no_sku"`
	ProductName      string             `json:"product_name"`
	Price            float64            `json:"price"`
	PromoPercentage  float64            `json:"promo_percentage"`
	IsPromo          int                `json:"flag_promo"`
	Description      string             `json:"description"`
	PictureUrl       string             `json:"picture_url"`
	Thumbnail        string             `json:"thumbnail"`
	StockOpname      int                `json:"stock_opname"`
	PriceInfo        string             `json:"price_info"`
	AccountType      string             `json:"account_type"`
	ListItemsPackage []ListItemsPackage `json:"list_items"`
}

type ListItemsPackage struct {
	Id          string  `json:"id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	PictureUrl  string  `json:"picture_url"`
	Thumbnail   string  `json:"thumbnail"`
	Price       float64 `json:"price"`
	Qty         int     `json:"qty"`
	SubTotal    float64 `json:"sub_total"`
}

func ToFindProductDesaByIdResponse(productDesa *entity.ProductsDesa, AccountType int, listItems []entity.ProductsPackageItems) (productDesaResponse FindProductDesaByIdResponse) {
	productDesaResponse.Id = productDesa.Id
	productDesaResponse.IdBrand = productDesa.ProductsMaster.IdBrand
	productDesaResponse.IdCategory = productDesa.ProductsMaster.IdCategory
	productDesaResponse.IdSubCategory = productDesa.ProductsMaster.IdSubCategory
	productDesaResponse.IdType = productDesa.IdType
	productDesaResponse.IdUnit = productDesa.ProductsMaster.IdUnit
	productDesaResponse.NoSku = productDesa.ProductsMaster.NoSku
	productDesaResponse.ProductName = productDesa.ProductsMaster.ProductName
	productDesaResponse.IsPromo = productDesa.IsPromo
	if AccountType == 1 {
		if productDesa.IsPromo == 1 {
			productDesaResponse.Price = productDesa.PricePromo
			productDesaResponse.PromoPercentage = productDesa.PercentagePromo
		} else {
			productDesaResponse.Price = productDesa.Price
		}
		productDesaResponse.AccountType = "User Biasa"
		productDesaResponse.PriceInfo = "Harga Normal"
	} else if AccountType == 2 {
		productDesaResponse.AccountType = "User Merchant"
		productDesaResponse.PriceInfo = "Harga Grosir"
		productDesaResponse.Price = productDesa.PriceGrosir
	}

	if productDesa.IdType == 1 {
		productDesaResponse.Description = productDesa.ProductsMaster.Description
		productDesaResponse.PictureUrl = productDesa.ProductsMaster.PictureUrl
		productDesaResponse.Thumbnail = productDesa.ProductsMaster.Thumbnail
		productDesaResponse.StockOpname = productDesa.StockOpname
	} else if productDesa.IdType == 2 {
		productDesaResponse.Description = productDesa.Description
		productDesaResponse.PictureUrl = productDesa.PictureUrl
		productDesaResponse.Thumbnail = productDesa.Thumbnail
		productDesaResponse.StockOpname = productDesa.StockOpname
	}

	if productDesa.IdType == 2 {
		listItemsResponses := []ListItemsPackage{}
		for _, listItem := range listItems {
			listItemResponse := ListItemsPackage{}
			listItemResponse.Id = listItem.Id
			listItemResponse.ProductName = listItem.ProductName
			productDesaResponse.Description = productDesa.ProductsMaster.Description
			productDesaResponse.PictureUrl = productDesa.ProductsMaster.PictureUrl
			productDesaResponse.Thumbnail = productDesa.ProductsMaster.Thumbnail
			listItemResponse.Price = listItem.Price
			listItemResponse.Qty = listItem.Qty
			listItemResponse.SubTotal = listItem.SubTotal
			listItemsResponses = append(listItemsResponses, listItemResponse)
		}
		productDesaResponse.ListItemsPackage = listItemsResponses
	}

	return productDesaResponse
}

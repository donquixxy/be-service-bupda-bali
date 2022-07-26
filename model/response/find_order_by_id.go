package response

import (
	"time"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type FindOrderByIdResponse struct {
	Id             string        `json:"id_order"`
	OrderName      string        `json:"order_name"`
	NumberOrder    string        `json:"number_order"`
	OrderStatus    int           `json:"order_status"`
	PaymentMethod  string        `json:"payment_method"`
	PaymentChannel string        `json:"payment_channel"`
	PaymentDueDate time.Time     `json:"payment_due_date"`
	SubTotal       float64       `json:"sub_total"`
	ShippingCost   float64       `json:"shipping_cost"`
	PaymentPoint   float64       `json:"payment_point"`
	PaymentFee     float64       `json:"payment_fee"`
	PaymentCash    float64       `json:"payment_cash"`
	TotalBill      float64       `json:"total_bill"`
	OrdersItems    []OrdersItems `json:"order_items"`
}

type OrdersItems struct {
	Id            string  `json:"id_item_order"`
	IdProductDesa string  `json:"id_product_desa"`
	Price         float64 `json:"price"`
	TotalPrice    float64 `json:"total_price"`
	ProductName   string  `json:"product_name"`
	Description   string  `json:"description"`
	PictureUrl    string  `json:"picture_url"`
	Thumbnail     string  `json:"thumbnail"`
	Qty           int     `json:"qty"`
	FlagPromo     int     `json:"flag_promo"`
}

func ToFindOrderByIdResponse(order *entity.Order, orderItems []entity.OrderItem) (orderResponse FindOrderByIdResponse) {
	orderResponse.Id = order.Id
	orderResponse.OrderName = order.OrderName
	orderResponse.NumberOrder = order.NumberOrder
	orderResponse.OrderStatus = order.OrderStatus
	orderResponse.PaymentMethod = order.PaymentMethod
	orderResponse.PaymentChannel = order.PaymentChannel
	orderResponse.PaymentDueDate = order.PaymentDueDate.Time
	orderResponse.SubTotal = order.SubTotal
	orderResponse.ShippingCost = order.ShippingCost
	orderResponse.PaymentPoint = order.PaymentPoint
	orderResponse.PaymentFee = order.PaymentFee
	orderResponse.PaymentCash = order.PaymentCash
	orderResponse.TotalBill = order.TotalBill

	var orderItemsResponses []OrdersItems
	for _, orderItem := range orderItems {
		var orderItemResponse OrdersItems
		orderItemResponse.Id = orderItem.Id
		orderItemResponse.IdProductDesa = orderItem.IdProductDesa
		if orderItem.FlagPromo == 1 {
			orderItemResponse.Price = orderItem.PriceAfterDiscount
		} else {
			orderItemResponse.Price = orderItem.Price
		}
		orderItemResponse.TotalPrice = orderItem.TotalPrice
		orderItemResponse.ProductName = orderItem.ProductName
		orderItemResponse.Description = orderItem.Description
		orderItemResponse.PictureUrl = orderItem.PictureUrl
		orderItemResponse.Thumbnail = orderItem.Thumbnail
		orderItemResponse.Qty = orderItem.Qty
		orderItemResponse.FlagPromo = orderItem.FlagPromo
		orderItemsResponses = append(orderItemsResponses, orderItemResponse)
	}

	orderResponse.OrdersItems = orderItemsResponses
	return orderResponse
}

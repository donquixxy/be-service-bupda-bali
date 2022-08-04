package response

import (
	"time"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type FindOrderByIdResponse struct {
	Id               string        `json:"id_order"`
	ProductType      string        `json:"product_type"`
	OrderType        int           `json:"order_type"`
	NumberOrder      string        `json:"number_order"`
	OrderStatus      int           `json:"order_status"`
	PaymentMethod    string        `json:"payment_method"`
	PaymentChannel   string        `json:"payment_channel"`
	PaymentDueDate   time.Time     `json:"payment_due_date"`
	SubTotal         float64       `json:"sub_total"`
	ShippingCost     float64       `json:"shipping_cost"`
	PaymentPoint     float64       `json:"payment_point"`
	PaymentFee       float64       `json:"payment_fee"`
	PaymentName      string        `json:"payment_name"`
	BankName         string        `json:"bank_name"`
	BankLogo         string        `json:"bank_logo"`
	PaymentNumber    string        `json:"payment_number"`
	PaymentCash      float64       `json:"payment_cash"`
	TotalBill        float64       `json:"total_bill"`
	AlamatPengiriman string        `json:"alamat_pengiriman"`
	CatatanKurir     string        `json:"catatan_kurir"`
	OrderDate        time.Time     `json:"order_date"`
	OrdersItems      []OrdersItems `json:"order_items"`
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

func ToFindOrderByIdResponse(order *entity.Order, orderItems []entity.OrderItem, payment *entity.PaymentChannel) (orderResponse FindOrderByIdResponse) {
	orderResponse.Id = order.Id
	orderResponse.ProductType = order.ProductType
	orderResponse.OrderType = order.OrderType
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
	orderResponse.AlamatPengiriman = order.AlamatPengiriman
	orderResponse.CatatanKurir = order.Catatan
	orderResponse.OrderDate = order.OrderedDate
	orderResponse.PaymentNumber = order.PaymentNo
	orderResponse.PaymentName = order.PaymentName
	orderResponse.BankName = payment.Name
	orderResponse.BankLogo = payment.Logo

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

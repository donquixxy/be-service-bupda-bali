package response

import (
	"time"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/ppob"
)

type FindOrderSembakoByIdResponse struct {
	Id               string               `json:"id_order"`
	ProductType      string               `json:"product_type"`
	OrderType        int                  `json:"order_type"`
	NumberOrder      string               `json:"number_order"`
	OrderStatus      int                  `json:"order_status"`
	PaymentMethod    string               `json:"payment_method"`
	PaymentChannel   string               `json:"payment_channel"`
	PaymentDueDate   time.Time            `json:"payment_due_date"`
	SubTotal         float64              `json:"sub_total"`
	ShippingCost     float64              `json:"shipping_cost"`
	PaymentPoint     float64              `json:"payment_point"`
	PaymentFee       float64              `json:"payment_fee"`
	PaymentName      string               `json:"payment_name"`
	BankName         string               `json:"bank_name"`
	BankLogo         string               `json:"bank_logo"`
	PaymentNumber    string               `json:"payment_number"`
	PaymentCash      float64              `json:"payment_cash"`
	TotalBill        float64              `json:"total_bill"`
	AlamatPengiriman string               `json:"alamat_pengiriman"`
	CatatanKurir     string               `json:"catatan_kurir"`
	OrderDate        time.Time            `json:"order_date"`
	OrdersItems      []OrdersItemsSembako `json:"order_items"`
}

type OrdersItemsSembako struct {
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

func ToFindOrderSembakoByIdResponse(order *entity.Order, orderItems []entity.OrderItem, payment *entity.PaymentChannel) (orderResponse FindOrderSembakoByIdResponse) {
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

	var orderItemsResponses []OrdersItemsSembako
	for _, orderItem := range orderItems {
		var orderItemResponse OrdersItemsSembako
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

type FindOrderPrepaidPulsaByIdResponse struct {
	Id                      string                  `json:"id_order"`
	ProductType             string                  `json:"product_type"`
	OrderType               int                     `json:"order_type"`
	NumberOrder             string                  `json:"number_order"`
	OrderStatus             int                     `json:"order_status"`
	PaymentMethod           string                  `json:"payment_method"`
	PaymentChannel          string                  `json:"payment_channel"`
	PaymentDueDate          time.Time               `json:"payment_due_date"`
	SubTotal                float64                 `json:"sub_total"`
	PaymentPoint            float64                 `json:"payment_point"`
	PaymentFee              float64                 `json:"payment_fee"`
	PaymentName             string                  `json:"payment_name"`
	BankName                string                  `json:"bank_name"`
	BankLogo                string                  `json:"bank_logo"`
	PaymentNumber           string                  `json:"payment_number"`
	PaymentCash             float64                 `json:"payment_cash"`
	TotalBill               float64                 `json:"total_bill"`
	OrderDate               time.Time               `json:"order_date"`
	OrdersItemsPrepaidPulsa OrdersItemsPrepaidPulsa `json:"order_items"`
}

type OrdersItemsPrepaidPulsa struct {
	TrId               int    `json:"tr_id"`
	RefId              string `json:"ref_id"`
	ProductCode        string `json:"product_code"`
	IconUrl            string `json:"icon_url"`
	ProductName        string `json:"product_name"`
	ProductDescription string `json:"product_description"`
	CustomerId         string `json:"customer_id"`
	Operator           string `json:"operator"`
	ActivePeriod       string `json:"active_periode"`
	StatusTopUp        int    `json:"status_topup"`
}

func ToFindOrderPrepaidPulsaByIdResponse(order *entity.Order, orderItemPpob *entity.OrderItemPpob, detailPpobPrepaidPulsa *entity.PpobDetailPrepaidPulsa, payment *entity.PaymentChannel) (orderResponse FindOrderPrepaidPulsaByIdResponse) {
	orderResponse.Id = order.Id
	orderResponse.ProductType = order.ProductType
	orderResponse.OrderType = order.OrderType
	orderResponse.NumberOrder = order.NumberOrder
	orderResponse.OrderStatus = order.OrderStatus
	orderResponse.PaymentMethod = order.PaymentMethod
	orderResponse.PaymentChannel = order.PaymentChannel
	orderResponse.PaymentDueDate = order.PaymentDueDate.Time
	orderResponse.SubTotal = orderItemPpob.TotalTagihan
	orderResponse.PaymentPoint = order.PaymentPoint
	orderResponse.PaymentFee = order.PaymentFee
	orderResponse.PaymentCash = order.PaymentCash
	orderResponse.TotalBill = order.TotalBill
	orderResponse.OrderDate = order.OrderedDate
	orderResponse.PaymentNumber = order.PaymentNo
	orderResponse.PaymentName = order.PaymentName
	orderResponse.BankName = payment.Name
	orderResponse.BankLogo = payment.Logo
	orderResponse.OrdersItemsPrepaidPulsa.TrId = orderItemPpob.TrId
	orderResponse.OrdersItemsPrepaidPulsa.RefId = orderItemPpob.RefId
	orderResponse.OrdersItemsPrepaidPulsa.ProductCode = orderItemPpob.ProductCode
	orderResponse.OrdersItemsPrepaidPulsa.IconUrl = orderItemPpob.IconUrl
	orderResponse.OrdersItemsPrepaidPulsa.ProductName = detailPpobPrepaidPulsa.ProductName
	orderResponse.OrdersItemsPrepaidPulsa.ProductDescription = detailPpobPrepaidPulsa.ProductDescription
	orderResponse.OrdersItemsPrepaidPulsa.CustomerId = detailPpobPrepaidPulsa.CustomerId
	orderResponse.OrdersItemsPrepaidPulsa.Operator = detailPpobPrepaidPulsa.Operator
	orderResponse.OrdersItemsPrepaidPulsa.ActivePeriod = detailPpobPrepaidPulsa.ActivePeriod
	orderResponse.OrdersItemsPrepaidPulsa.StatusTopUp = detailPpobPrepaidPulsa.StatusTopUp
	return orderResponse
}

type FindOrderPrepaidPlnByIdResponse struct {
	Id                    string                `json:"id_order"`
	ProductType           string                `json:"product_type"`
	OrderType             int                   `json:"order_type"`
	NumberOrder           string                `json:"number_order"`
	OrderStatus           int                   `json:"order_status"`
	PaymentMethod         string                `json:"payment_method"`
	PaymentChannel        string                `json:"payment_channel"`
	PaymentDueDate        time.Time             `json:"payment_due_date"`
	SubTotal              float64               `json:"sub_total"`
	PaymentPoint          float64               `json:"payment_point"`
	PaymentFee            float64               `json:"payment_fee"`
	PaymentName           string                `json:"payment_name"`
	BankName              string                `json:"bank_name"`
	BankLogo              string                `json:"bank_logo"`
	PaymentNumber         string                `json:"payment_number"`
	PaymentCash           float64               `json:"payment_cash"`
	TotalBill             float64               `json:"total_bill"`
	OrderDate             time.Time             `json:"order_date"`
	OrdersItemsPrepaidPln OrdersItemsPrepaidPln `json:"order_items"`
}

type OrdersItemsPrepaidPln struct {
	TrId               int    `json:"tr_id"`
	RefId              string `json:"ref_id"`
	ProductCode        string `json:"product_code"`
	ProductName        string `json:"product_name"`
	ProductDescription string `json:"product_description"`
	CustomerId         string `json:"customer_id"`
	MeterNo            string `json:"meter_no"`
	SubscriberId       string `json:"subscriber_id"`
	CustomerName       string `json:"customer_name"`
	SegmentPower       string `json:"segment_power"`
	StatusTopUp        int    `json:"status_topup"`
	NoToken            string `json:"no_token"`
}

func ToFindOrderPrepaidPlnByIdResponse(order *entity.Order, orderItemPpob *entity.OrderItemPpob, detailPpobPrepaidPln *entity.PpobDetailPrepaidPln, payment *entity.PaymentChannel) (orderResponse FindOrderPrepaidPlnByIdResponse) {
	orderResponse.Id = order.Id
	orderResponse.ProductType = order.ProductType
	orderResponse.OrderType = order.OrderType
	orderResponse.NumberOrder = order.NumberOrder
	orderResponse.OrderStatus = order.OrderStatus
	orderResponse.PaymentMethod = order.PaymentMethod
	orderResponse.PaymentChannel = order.PaymentChannel
	orderResponse.PaymentDueDate = order.PaymentDueDate.Time
	orderResponse.SubTotal = orderItemPpob.TotalTagihan
	orderResponse.PaymentPoint = order.PaymentPoint
	orderResponse.PaymentFee = order.PaymentFee
	orderResponse.PaymentCash = order.PaymentCash
	orderResponse.TotalBill = order.TotalBill
	orderResponse.OrderDate = order.OrderedDate
	orderResponse.PaymentNumber = order.PaymentNo
	orderResponse.PaymentName = order.PaymentName
	orderResponse.BankName = payment.Name
	orderResponse.BankLogo = payment.Logo
	orderResponse.OrdersItemsPrepaidPln.TrId = orderItemPpob.TrId
	orderResponse.OrdersItemsPrepaidPln.RefId = orderItemPpob.RefId
	orderResponse.OrdersItemsPrepaidPln.ProductCode = orderItemPpob.ProductCode
	orderResponse.OrdersItemsPrepaidPln.ProductName = detailPpobPrepaidPln.ProductName
	orderResponse.OrdersItemsPrepaidPln.ProductDescription = detailPpobPrepaidPln.ProductDescription
	orderResponse.OrdersItemsPrepaidPln.CustomerId = detailPpobPrepaidPln.CustomerId
	orderResponse.OrdersItemsPrepaidPln.MeterNo = detailPpobPrepaidPln.MeterNo
	orderResponse.OrdersItemsPrepaidPln.SubscriberId = detailPpobPrepaidPln.SubscriberId
	orderResponse.OrdersItemsPrepaidPln.CustomerName = detailPpobPrepaidPln.CustomerName
	orderResponse.OrdersItemsPrepaidPln.SegmentPower = detailPpobPrepaidPln.SegmentPower
	orderResponse.OrdersItemsPrepaidPln.StatusTopUp = detailPpobPrepaidPln.StatusTopUp
	orderResponse.OrdersItemsPrepaidPln.NoToken = detailPpobPrepaidPln.NoToken
	return orderResponse
}

type FindOrderPostpaidPlnByIdResponse struct {
	Id                     string                 `json:"id_order"`
	ProductType            string                 `json:"product_type"`
	OrderType              int                    `json:"order_type"`
	NumberOrder            string                 `json:"number_order"`
	OrderStatus            int                    `json:"order_status"`
	PaymentMethod          string                 `json:"payment_method"`
	PaymentChannel         string                 `json:"payment_channel"`
	PaymentDueDate         time.Time              `json:"payment_due_date"`
	SubTotal               float64                `json:"sub_total"`
	PaymentPoint           float64                `json:"payment_point"`
	PaymentFee             float64                `json:"payment_fee"`
	PaymentName            string                 `json:"payment_name"`
	BankName               string                 `json:"bank_name"`
	BankLogo               string                 `json:"bank_logo"`
	PaymentNumber          string                 `json:"payment_number"`
	PaymentCash            float64                `json:"payment_cash"`
	TotalBill              float64                `json:"total_bill"`
	OrderDate              time.Time              `json:"order_date"`
	OrdersItemsPostpaidPln OrdersItemsPostpaidPln `json:"order_items"`
}

type OrdersItemsPostpaidPln struct {
	TrId              int                 `json:"tr_id"`
	RefId             string              `json:"ref_id"`
	CustomerId        string              `json:"customer_id"`
	Tarif             string              `json:"tarif"`
	Daya              int                 `json:"daya"`
	LembarTagihan     string              `json:"lembar_tagihan"`
	Period            string              `json:"period"`
	StatusTopUp       int                 `json:"status_topup"`
	PostpaidPlnDetail []PostpaidPlnDetail `json:"detail_tagihan"`
}

type PostpaidPlnDetail struct {
	Periode     string  `json:"periode"`
	NilaiTgihan string  `json:"nilai_tagihan"`
	Admin       string  `json:"admin"`
	Denda       string  `json:"denda"`
	Total       float64 `json:"total"`
}

func ToFindOrderPostpaidPlnByIdResponse(order *entity.Order, orderItemPpob *entity.OrderItemPpob, detailPpobPostpaidPln *entity.PpobDetailPostpaidPln, payment *entity.PaymentChannel, detailTagihans []ppob.InquiryPostpaidPlnDetail) (orderResponse FindOrderPostpaidPlnByIdResponse) {
	orderResponse.Id = order.Id
	orderResponse.ProductType = order.ProductType
	orderResponse.OrderType = order.OrderType
	orderResponse.NumberOrder = order.NumberOrder
	orderResponse.OrderStatus = order.OrderStatus
	orderResponse.PaymentMethod = order.PaymentMethod
	orderResponse.PaymentChannel = order.PaymentChannel
	orderResponse.PaymentDueDate = order.PaymentDueDate.Time
	orderResponse.SubTotal = orderItemPpob.TotalTagihan
	orderResponse.PaymentPoint = order.PaymentPoint
	orderResponse.PaymentFee = order.PaymentFee
	orderResponse.PaymentCash = order.PaymentCash
	orderResponse.TotalBill = order.TotalBill
	orderResponse.OrderDate = order.OrderedDate
	orderResponse.PaymentNumber = order.PaymentNo
	orderResponse.PaymentName = order.PaymentName
	orderResponse.BankName = payment.Name
	orderResponse.BankLogo = payment.Logo
	orderResponse.OrdersItemsPostpaidPln.CustomerId = detailPpobPostpaidPln.CustomerId
	orderResponse.OrdersItemsPostpaidPln.RefId = detailPpobPostpaidPln.RefId
	orderResponse.OrdersItemsPostpaidPln.Tarif = detailPpobPostpaidPln.Tarif
	orderResponse.OrdersItemsPostpaidPln.Daya = detailPpobPostpaidPln.Daya
	orderResponse.OrdersItemsPostpaidPln.Period = detailPpobPostpaidPln.Period
	orderResponse.OrdersItemsPostpaidPln.LembarTagihan = detailPpobPostpaidPln.LembarTagihan
	orderResponse.OrdersItemsPostpaidPln.StatusTopUp = detailPpobPostpaidPln.StatusTopUp

	var postpaidPlnDetails []PostpaidPlnDetail
	for _, detailTagihan := range detailTagihans {
		var postpaidPlnDetail PostpaidPlnDetail
		postpaidPlnDetail.Admin = detailTagihan.Admin
		postpaidPlnDetail.NilaiTgihan = detailTagihan.NilaiTgihan
		postpaidPlnDetail.Periode = detailTagihan.Periode
		postpaidPlnDetail.Denda = detailTagihan.Denda
		postpaidPlnDetail.Total = detailTagihan.Total
		postpaidPlnDetails = append(postpaidPlnDetails, postpaidPlnDetail)
	}
	orderResponse.OrdersItemsPostpaidPln.PostpaidPlnDetail = postpaidPlnDetails

	return orderResponse
}

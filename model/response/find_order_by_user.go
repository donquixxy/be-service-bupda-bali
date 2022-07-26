package response

import (
	"time"

	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type FindOrderByUserResponse struct {
	Id          string    `json:"id_order"`
	IdUser      string    `json:"id_user"`
	OrderName   string    `json:"order_name"`
	OrderStatus int       `json:"order_status"`
	OrderedDate time.Time `json:"order_date"`
	TotalBill   float64   `json:"total_bill"`
}

func ToFindOrderByUserResponse(orders []entity.Order) (orderResponses []FindOrderByUserResponse) {
	for _, order := range orders {
		var orderResponse FindOrderByUserResponse
		orderResponse.Id = order.Id
		orderResponse.IdUser = order.IdUser
		orderResponse.OrderName = order.OrderName
		orderResponse.OrderStatus = order.OrderStatus
		orderResponse.OrderedDate = order.OrderedDate
		orderResponse.TotalBill = order.TotalBill
		orderResponses = append(orderResponses, orderResponse)
	}
	return orderResponses
}

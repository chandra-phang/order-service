package v1

import (
	"order-service/lib"
	"order-service/model"
)

type ListOrdersDTO struct {
	Orders []OrderDTO `json:"orders"`
}

type OrderDTO struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (dto *ListOrdersDTO) ConvertFromOrdersEntity(orders []model.Order) *ListOrdersDTO {
	resp := &ListOrdersDTO{}
	for _, order := range orders {
		orderDTO := OrderDTO{
			ID:        order.ID,
			UserID:    order.UserID,
			Status:    string(order.Status),
			CreatedAt: lib.ConvertToEpoch(order.CreatedAt),
			UpdatedAt: lib.ConvertToEpoch(order.UpdatedAt),
		}
		resp.Orders = append(resp.Orders, orderDTO)
	}

	return resp
}

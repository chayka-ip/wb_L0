package store

import models "level_zero/internal/app/model"

type OrderRepository interface {
	Create(*models.Order) error
	GetByOrderUid(string) (*models.Order, error)
	ListOrders(int) ([]*models.Order, error)
}

package teststore

import (
	models "level_zero/internal/app/model"
	"level_zero/store"
)

type Store struct {
	orderRepository *OrderRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) Order() store.OrderRepository {
	if s.orderRepository == nil {
		s.orderRepository = &OrderRepository{
			Store:  s,
			orders: make(map[string]*models.Order),
		}
	}
	return s.orderRepository
}

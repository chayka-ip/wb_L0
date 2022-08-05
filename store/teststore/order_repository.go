package teststore

import (
	"encoding/json"
	"errors"
	models "level_zero/internal/app/model"
	"level_zero/order"
	"level_zero/store"
)

type OrderRepository struct {
	Store  *Store
	orders map[string]*models.Order
}

func (r *OrderRepository) Create(ord *models.Order) error {
	if _, ok := r.orders[ord.OrderUID]; ok {
		return errors.New("Item already exists")
	}

	vOrder := &order.Order{}
	err := json.Unmarshal(ord.Data, vOrder)
	if err != nil {
		return err
	}

	r.orders[ord.OrderUID] = ord

	return nil
}

func (r *OrderRepository) GetByOrderUid(orderUID string) (*models.Order, error) {
	ord, ok := r.orders[orderUID]
	if ok {
		return ord, nil
	}
	return nil, store.ErrRecordNotFound
}

func (r *OrderRepository) ListOrders(count int) ([]*models.Order, error) {
	l := len(r.orders)
	if count > l {
		count = l
	}

	var out []*models.Order
	for i := range r.orders {
		if len(out) >= count {
			break
		}
		out = append(out, r.orders[i])
	}
	return out, nil
}

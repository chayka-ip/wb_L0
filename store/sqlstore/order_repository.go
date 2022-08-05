package sqlstore

import (
	models "level_zero/internal/app/model"
	"level_zero/store"
)

type OrderRepository struct {
	store *Store
}

var (
	queryCreate       = "INSERT INTO orders (order_uid, data) VALUES ($1, $2) RETURNING order_uid"
	queryGetOrderById = "SELECT order_uid, data FROM orders WHERE order_uid=$1"
	queryListOrders   = "SELECT order_uid, data FROM orders LIMIT $1"
)

func (r *OrderRepository) Create(order *models.Order) error {
	return r.store.db.QueryRow(queryCreate, order.OrderUID, order.Data).Scan(&order.OrderUID)
}

func (r *OrderRepository) GetByOrderUid(order_uid string) (*models.Order, error) {
	order := &models.Order{}
	err := r.store.db.QueryRow(queryGetOrderById, order_uid).Scan(&order.OrderUID, &order.Data)
	if err != nil {
		if err == store.ErrRecordNotFound {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return order, nil
}

func (r *OrderRepository) ListOrders(count int) ([]*models.Order, error) {
	rows, err := r.store.db.Query(queryListOrders, count)
	if err != nil {
		return nil, err
	}
	var orderList []*models.Order
	for rows.Next() {
		order := &models.Order{}
		if err := rows.Scan(&order.OrderUID, &order.Data); err != nil {
			return nil, err
		}
		orderList = append(orderList, order)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	return orderList, nil
}

package sqlstore

import (
	"database/sql"
	"level_zero/store"

	_ "github.com/lib/pq"
)

type Store struct {
	db              *sql.DB
	orderRepository *OrderRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Order() store.OrderRepository {
	if s.orderRepository == nil {
		s.orderRepository = &OrderRepository{
			store: s,
		}
	}
	return s.orderRepository
}

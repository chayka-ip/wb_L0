package sqlstore_test

import (
	models "level_zero/internal/app/model"
	"level_zero/store"
	"level_zero/store/sqlstore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseUrl)
	defer teardown(orders_table)
	s := sqlstore.New(db)

	newOrder := models.TestOrder()
	err := s.Order().Create(newOrder)

	// can create order
	assert.NoError(t, err)

	// can't add duplicate
	err = s.Order().Create(newOrder)
	assert.Error(t, err)

	// can't add order with invalid data
	newOrder.OrderUID += "1"
	newOrder.Data = append(newOrder.Data, 123)
	err = s.Order().Create(newOrder)
	assert.Error(t, err)

}

func TestOrderRepository_GetByOrderUid(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseUrl)
	defer teardown(orders_table)
	s := sqlstore.New(db)

	ord := models.TestOrder()

	_, err := s.Order().GetByOrderUid(ord.OrderUID)
	assert.Error(t, store.ErrRecordNotFound)

	s.Order().Create(ord)
	fromDB, err := s.Order().GetByOrderUid(ord.OrderUID)
	assert.NoError(t, err)
	assert.NotNil(t, fromDB)
	assert.Equal(t, fromDB.OrderUID, ord.OrderUID)
}

func TestOrderRepository_ListOrders(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseUrl)
	defer teardown(orders_table)
	s := sqlstore.New(db)

	const totalItemsDB = 10
	var orderList [totalItemsDB]*models.Order
	for i := 0; i < totalItemsDB; i++ {
		orderList[i] = models.TestOrder()
	}

	numListItems := 5

	res1, err := s.Order().ListOrders(numListItems)
	assert.NoError(t, err)
	assert.Empty(t, res1)

	for i := range orderList {
		err = s.Order().Create(orderList[i])
		assert.NoError(t, err)
	}

	res2, err := s.Order().ListOrders(numListItems)
	assert.NoError(t, err)
	assert.Equal(t, len(res2), numListItems)
}

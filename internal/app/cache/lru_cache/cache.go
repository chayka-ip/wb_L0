package lru_cache

import (
	"encoding/json"
	"level_zero/internal/app/logger"
	"level_zero/store"

	lru "github.com/hashicorp/golang-lru"
)

type Cache struct {
	cache  *lru.Cache // [order_uid] : models.Order.Data
	size   int
	store  store.Store
	logger logger.Logger
}

func NewCache(size int, store store.Store, logger logger.Logger) (*Cache, error) {
	c, err := lru.New(size)
	if err != nil {
		return nil, err
	}

	cache := &Cache{
		cache:  c,
		size:   size,
		store:  store,
		logger: logger,
	}

	cache.Recover()
	return cache, nil
}

func (c *Cache) AddToStore(key string, value json.RawMessage) {
	c.cache.ContainsOrAdd(key, value)
}

// returns models.Order.Data or error
func (c *Cache) Get(key string) (json.RawMessage, error) {
	value, ok := c.cache.Get(key)
	if ok {
		v, _ := value.(json.RawMessage)
		return v, nil
	}

	order, err := c.store.Order().GetByOrderUid(key)
	if err != nil {
		return nil, err
	}

	c.AddToStore(key, order.Data)
	return order.Data, nil
}

func (c *Cache) Recover() {
	orderList, err := c.store.Order().ListOrders(c.size)
	if err != nil {
		c.logger.Info("Failed to restore cache from database")
		return
	}

	for _, data := range orderList {
		c.AddToStore(data.OrderUID, data.Data)
	}
	c.logger.Info("Cache is recovered from database")
}

package models

import (
	"level_zero/order"
	"log"

	"github.com/bxcodec/faker"
)

func TestOrder() *Order {
	ord := &order.Order{}
	err := faker.FakeData(ord)
	if err != nil {
		log.Fatal(err)
	}
	return MakeOrderModelFromStruct(ord)
}

func TestOrderBadData() *Order {
	ord := TestOrder()
	ord.Data = append(ord.Data, 123)
	return ord
}

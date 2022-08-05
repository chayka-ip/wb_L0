package models

import (
	"encoding/json"
	"level_zero/order"
	"log"
)

type Order struct {
	OrderUID string          `json:"order_uid"`
	Data     json.RawMessage `json:"data"`
}

func MakeOrderModelFromStruct(order *order.Order) *Order {
	data, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}
	return &Order{
		OrderUID: order.OrderUID,
		Data:     data,
	}
}

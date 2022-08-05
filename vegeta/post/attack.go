package main

import (
	"fmt"
	"level_zero/order"
	vegeta_utils "level_zero/vegeta/utils"
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	attackName := "Attack: POST ORDER"
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 1000 * time.Millisecond

	targeter := getNewTargeter()
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	var successCount int64 = 0
	for res := range attacker.Attack(targeter, rate, duration, attackName) {
		metrics.Add(res)
		if res.Code == http.StatusCreated {
			successCount++
		}
	}
	vegeta_utils.PrintMetrics(attackName, &metrics, successCount, rate.String())
	defer metrics.Close()
}

func getNewTargeter() vegeta.Targeter {
	return func(tg *vegeta.Target) error {
		if tg == nil {
			return vegeta.ErrNilTarget
		}
		socket := "http://127.0.0.1:8080"
		targetAPI := "api/add_order"
		URL := fmt.Sprintf("%s/%s", socket, targetAPI)

		tg.Method = http.MethodPost
		tg.URL = URL
		tg.Header = http.Header{
			"Content-Type": {"application/json"},
		}
		tg.Body = order.GetRandomOrderEncodedJson()
		return nil
	}
}

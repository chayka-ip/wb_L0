package main

import (
	"fmt"
	vegeta_utils "level_zero/vegeta/utils"
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// TODO: make function to create new order before attack

func main() {
	attackName := "GET ORDER BY UID"
	rate := vegeta.Rate{Freq: 40000, Per: time.Second}
	duration := 600 * time.Millisecond
	socket := "http://127.0.0.1:8080"
	targetAPI := "api/order"
	orderUID := "test12"

	URL := fmt.Sprintf("%s/%s?id=%s", socket, targetAPI, orderUID)
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: http.MethodGet,
		URL:    URL,
	})
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	var successCount int64 = 0
	for res := range attacker.Attack(targeter, rate, duration, attackName) {
		metrics.Add(res)
		if res.Code == http.StatusOK {
			successCount++
		}
	}
	fmt.Println(metrics.Latencies.P90, metrics.Latencies.Total)
	vegeta_utils.PrintMetrics(attackName, &metrics, successCount, rate.String())
	metrics.Close()
}

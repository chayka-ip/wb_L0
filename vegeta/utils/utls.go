package vegeta_utils

import (
	"fmt"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func PrintMetrics(attackName string, metrics *vegeta.Metrics, successfulCount int64, sendRate string) {
	successRate := int(float64(successfulCount) * 100.0 / float64(metrics.Requests))
	p := "================================\n"
	s := p
	s += fmt.Sprintf("[ATTACK]: %s\n", attackName)
	s += fmt.Sprintf("[REQUESTS] : %d\n", metrics.Requests)
	// s += fmt.Sprintf("[LATENCIES]\n")
	// s += fmt.Sprintf("    90th percentile : %d\n", metrics.Latencies.P90)
	// s += fmt.Sprintf("    mean : %d\n", metrics.Latencies.Mean)
	s += fmt.Sprintf("[SEND RATE] %s\n", sendRate)
	s += fmt.Sprintf("[SUCCESS %%] %d\n", successRate)
	s += p
	fmt.Println(s)
}

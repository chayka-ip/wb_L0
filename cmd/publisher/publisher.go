package main

import (
	"encoding/json"
	"fmt"
	"level_zero/internal/app/apiserver"
	"level_zero/order"
	"log"
	"math/rand"
	"time"

	"github.com/bxcodec/faker"
	"github.com/nats-io/stan.go"
)

func main() {
	config := apiserver.LoadConfig()
	clusterId := config.NatsClusterId
	pubId := config.NatsPubliserId
	failConnMsg := fmt.Sprintf(
		`Could not connect publisher to STAN	with settings: 
		 cluster: %v, pub_id: %v \n`, clusterId, pubId)

	conn, err := stan.Connect(clusterId, pubId)
	if err != nil {
		log.Fatal(failConnMsg)
	}

	msgPerSec := config.PublisherSendRate
	if msgPerSec < 1 {
		msgPerSec = 1
	}

	workTime := int64(config.PublisherWorkTime * float64(time.Second))
	timeStop := time.Now().Add(time.Duration(workTime)).UnixNano()

	badDataChance := config.BadDataChance

	sendRate := float64(time.Second) / float64(msgPerSec)
	waitDuration := time.Duration(sendRate)

	numSends := 0

	for {
		if time.Now().UnixNano() >= timeStop {
			break
		}
		broken := shouldBreakData(badDataChance)
		data := getMockData(broken)
		conn.PublishAsync(clusterId, data, AckHandler)
		numSends++

		s := fmt.Sprintf("%d STAN send iteration. ", numSends)
		if broken {
			s += "Sending brokend data..."
		}
		s += "\n"

		fmt.Println(s)
		time.Sleep(waitDuration)
	}

}

func shouldBreakData(chance float32) bool {
	if (chance < 0.0) || (chance > 1.0) {
		chance = 0.0
	}
	return chance >= rand.Float32()
}

func getMockData(bBroken bool) []byte {
	ord := order.Order{}
	faker.FakeData(&ord)
	out, err := json.Marshal(ord)
	if err != nil {
		log.Fatal(err)
	}

	if bBroken {
		out = append(out, 123)
	}

	return out
}

func AckHandler(string, error) {

}

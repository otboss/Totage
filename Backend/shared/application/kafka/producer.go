package kafka

import (
	"log"
	"waterordersystem/GOBACKEND/utils"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var Producer kafka.Producer

func init() {
	Producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": utils.Env.KAFKA_BOOTSTRAP_SERVERS,
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
	}
	defer Producer.Close()

}

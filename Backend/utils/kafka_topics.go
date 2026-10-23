package utils

type kafkaTopics struct {
	OrderEvent string
}

var KafkaTopics kafkaTopics = kafkaTopics{
	OrderEvent: "order_event",
}

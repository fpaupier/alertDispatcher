package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
)

const (
	bootstrapServers = "pkc-4r297.europe-west1.gcp.confluent.cloud:9092"
	ccloudAPIKey     = ConfluentApiKey
	ccloudAPISecret  = ConfluentSecret
)

func publish(msg []byte) {
	topic := "alert-topic"
	// Produce a new record to the topic...
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname: %v\n", err)
	}
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"client.id":         hostname,
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_SSL",
		"sasl.username":     ccloudAPIKey,
		"sasl.password":     ccloudAPISecret})

	if err != nil {
		panic(fmt.Sprintf("Failed to create producer: %s", err))
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg},
		nil)

	if err != nil {
		panic(fmt.Sprintf("failed to produce message : %b, got error: %s", msg, err))
	}

	// Wait for delivery report
	e := <-producer.Events()

	message := e.(*kafka.Message)
	if message.TopicPartition.Error != nil {
		fmt.Printf("failed to deliver message: %v\n",
			message.TopicPartition)
	} else {
		fmt.Printf("delivered to topic %s [%d] at offset %v\n",
			*message.TopicPartition.Topic,
			message.TopicPartition.Partition,
			message.TopicPartition.Offset)
	}

	producer.Close()

}

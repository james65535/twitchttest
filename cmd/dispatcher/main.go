package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

var (
	kafkaAddress string
	topic string
	groupId string
)

// Handles writing to Kafka broker
func kafkaWrite(topic, msg string) error {
	partition := 0
	conn, connectErr := kafka.DialLeader(context.Background(), "tcp", kafkaAddress, topic, partition)
	if connectErr != nil {
		return connectErr
	}
	conn.SetWriteDeadline(time.Now().Add(10*time.Second))
	_, writeErr := conn.WriteMessages(kafka.Message{Value: []byte(msg)})
	if writeErr != nil {
		return writeErr
	}
	conn.Close()
	return nil
}

func init() {
	// Kafka setup
	if os.Getenv("KAFKAADDRESS") != "" {
		kafkaAddress = os.Getenv("KAFKAADDRESS")
		log.Printf("Using kafka server: %v\n", kafkaAddress)
	} else {
		kafkaAddress = "127.0.0.1:9092"
		log.Printf("No kafka server address specified, defaulting to: %v\n", kafkaAddress)
	}

	if os.Getenv("TOPIC") != "" {
		topic = os.Getenv("TOPIC")
		log.Printf("Using kafka topic: %v\n", topic)
	} else {
		topic = "test-topic"
		log.Printf("No Kafka Topic specified, defaulting to: %v\n", topic)
	}

	if os.Getenv("GROUPID") != "" {
		groupId = os.Getenv("GROUPID")
		log.Printf("Using kafka Group ID: %v\n", groupId)
	} else {
		groupId = "test-dispatch-group"
		log.Printf("No kafka Group ID specified, defaulting to: %v\n", groupId)
	}
}


func main() {
	log.Printf("Starting dispatcher\n")
	if topic == "test-topic" {
		log.Printf("Running topic write test\n")
		// write to test topic
		msg := fmt.Sprintf("test: %s", time.Now())
		writeErr := kafkaWrite(topic, msg)
		if writeErr != nil {
			log.Printf("Error writing message: %s\n", writeErr)
		} else {
			log.Printf("Test write completed\n")
		}
	}

	// TODO check dispatch queue for new messages and process

}

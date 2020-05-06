package main

import (
	"context"
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

type Data struct {
	Followed_at string
	From_id string
	From_name string
	To_id string
	To_name string
}

type UserMsg struct {
	Header map[string][]string
	Body []Data
}

// Handles writing to Kafka broker
func kafkaWrite(topic string, msg []byte) error {
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		kafkaAddress,
		topic,
		0)
	if err != nil {
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(10*time.Second))
	if _, err := conn.WriteMessages(
		kafka.Message{
			Value: []byte(msg)})
	err != nil {
		return err
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
		//header := map[string]string{"link":`"[<https://api.twitch.tv/helix/webhooks/hub>; rel="hub", <https://api.twitch.tv/helix/users/follows?first=1&to_id=188951100>; rel="self"]"`}
		msg1 := UserMsg{
			map[string]string{"link":`"[<https://api.twitch.tv/helix/webhooks/hub>; rel="hub", <https://api.twitch.tv/helix/users/follows?first=1&to_id=188951100>; rel="self"]"`},
			Body: []Data{},
		}
		msg := []byte("{\"header\": \"[<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users/follows?first=1&to_id=188951100>; rel=\"self\"]\",\"body\": {\"data\": [{\"followed_at\": \"2020-04-22T22:26:18Z\",\"from_id\": \"59480475\",\"from_name\": \"eventualdecline\",\"to_id\": \"188951100\",\"to_name\": \"james65535\"}]}}")

		//msg := fmt.Sprintf("test: %s", time.Now())
		if err := kafkaWrite(topic, msg)
		err != nil {
			log.Printf("Error writing message: %s\n", err)
		}

		log.Printf("Test write completed\n")
	}

	// TODO check dispatch queue for new messages and process

}

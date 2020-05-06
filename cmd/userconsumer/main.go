package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
)

var (
	kafkaAddress string
	topic string
	groupId string
)
type data struct {
	Followed_at string
	From_id string
	From_name string
	To_id string
	To_name string
}

type userMsg struct {
	Header map[string][]string
	Body []data
}

func kafkaReader(topic, gid, addr string) error {
	r := kafka.NewReader(
		kafka.ReaderConfig{
		Brokers: []string{addr},
		GroupID: gid,
		Topic: topic,
		MinBytes: 10e3,
		MaxBytes: 10e6})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		log.Printf("Message at topic/partition/offset %v/%v/%v: %s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		return nil
	}
}

/*
// TODO remove in favour of reader
// Handles reading from Kafka broker
func kafkaRead(topic string) error {
	partition := 0
	conn,err := kafka.DialLeader(context.Background(), "tcp", kafkaAddress, topic, partition)
	if err != nil {
		return err
	}

	conn.SetReadDeadline(time.Now().Add(10*time.Second))
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		if _, err := batch.Read(b); err != nil {
			break
		}
		log.Printf("kafka msg body: %v\n", string(b))
	}

	batch.Close()
	conn.Close()
	return nil
}
*/

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
		groupId = "test-group"
		log.Printf("No kafka Group ID specified, defaulting to: %v\n", groupId)
	}
}

func main(){
	log.Printf("Starting consumer\n")
	// readErr := kafkaRead(topic)
	if topic == "test-topic" {
		log.Printf("Running topic read test\n")
		// write to test topic
		if err := kafkaReader(
			topic,
			groupId,
			kafkaAddress)
		err != nil {
			log.Printf("Error reading message: %s\n", err)
		}

		log.Printf("Test read completed\n")
	}
}
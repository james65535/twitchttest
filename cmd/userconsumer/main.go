package main

import (
	"github.com/segmentio/kafka-go"
	"log"
	"time"
	"context"
	"os"
)

var (
	kafkaAddress string
)


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
		_, err := batch.Read(b)
		if err != nil {
			break
		}
		log.Printf("kafka msg body: %v\n", string(b))
	}

	batch.Close()
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
}

func main(){
	readErr := kafkaRead("twitch")
	if readErr != nil {
		log.Fatalf("Cannot read: %s\n", readErr)
	}
	log.Printf("Done kafka read\n")
}
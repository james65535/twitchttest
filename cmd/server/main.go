package main

import (
	"context"
	"fmt"
	kafka "github.com/segmentio/kafka-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	address string
	kafkaAddress string
)

func kafkaWrite() error {
	// to produce messages
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaAddress, topic, partition)
    if err != nil {
    	return err
	}
	conn.SetWriteDeadline(time.Now().Add(10*time.Second))
	conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	conn.Close()
	return nil
}

func kafkaRead() error {
	// to consume messages
	topic := "my-topic"
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
		fmt.Println(string(b))
	}

	batch.Close()
	conn.Close()
	return nil
}

func init() {
	// Configure GraphDB Access
	if os.Getenv("ADDRESS") != "" {
		address = os.Getenv("ADDRESS")
	} else {
		log.Printf("No server address specified\n")
		address = ":8080"
	}
	if os.Getenv("KAFKAADDRESS") != "" {
		kafkaAddress = os.Getenv("KAFKAADDRESS")
	} else {
		log.Printf("No kafka server address specified\n")
		kafkaAddress = "127.0.0.1:9092"
	}
}

func main() {
	log.Printf("Server starting")
	submit := func(w http.ResponseWriter, r *http.Request){
		if r.Method == "POST" {
			log.Printf("POST Response received!\n")
			log.Printf("Header: %s",r.Header)
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("Cannot read body: %s", err)
			}
			log.Printf("Body: %s",body)
		} else { return }
	}
	home := func(w http.ResponseWriter, r *http.Request){
		if r.Method == "GET" {
			log.Printf("GET Response received!\n")
			err := kafkaWrite()
			if err != nil {
				log.Printf("Kafka Write Error %v\n", err)
			}
			log.Printf("Done with Kafka Write\n")

			readErr := kafkaRead()
			if readErr != nil {
				log.Printf("Kafka Read Error %v\n", readErr)
			}
		} else {return}
	}
	http.HandleFunc("/submit", submit)
	http.HandleFunc("/", home)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("Error %v\n", err)
	}
}

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
	clientId string
	clientSecret string
	callbackUrl string
)

type notificationMsg struct {
	link string
	body string
}
func writeNotification() error {


	return nil
}

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

// Respond to callback challenge
func callback(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Receive subscription verify/deny challenges
		log.Printf("Callback GET request received from %v!\n", r.RemoteAddr)
		// TODO verify that subscription challenge is a verify or a deny. Check request payload

		query := r.URL.Query()
		log.Printf("URL Params: %s", query) // URL Params: map[hub.challenge:[GJvOdj0BQOXdjN5gP7giEDU0I_dUA-nzyOANT8Sp] hub.lease_seconds:[6400] hub.mode:[subscribe] hub.topic:[https://api.twitch.tv/helix/users/follows?first=1&to_id=188951100]]
		challenge := query.Get("hub.challenge")
		log.Printf("hub.challenge: %v\n", challenge)

		// Respond to subscription challenge verify request
		w.WriteHeader(200)
		w.Write([]byte(challenge))
	} else if r.Method == "POST" {
		// Receive subscription notification
		// TODO verify header "Twitch-Notification-Id" with secret hash
		log.Printf("Callback POST request received!\n")
		log.Printf("Callback Header: %s", r.Header)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Cannot read body: %s", err)
		}
		log.Printf("Callback Body: %s", body)
	}
}

func submit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Printf("Submit POST request received!\n")
		log.Printf("Header: %s", r.Header)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Cannot read body: %s", err)
		}
		log.Printf("Body: %s", body)
	} else if r.Method == "GET" {
		log.Printf("Submit GET request received!\n")
		log.Printf("Header: %s", r.Header)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Cannot read body: %s", err)
		}
		log.Printf("Body: %s", body)
	}
}

func home(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		log.Printf("Home GET request received!\n")
		/*
			err := kafkaWrite()
			if err != nil {
				log.Printf("Kafka Write Error %v\n", err)
			}
			log.Printf("Done with Kafka Write\n")

			readErr := kafkaRead()
			if readErr != nil {
				log.Printf("Kafka Read Error %v\n", readErr)
			}
		*/
	} else {return}
}

func init() {
	// Webserver
	if os.Getenv("ADDRESS") != "" {
		address = os.Getenv("ADDRESS")
	} else {
		address = ":8080"
		log.Printf("No server address specified, defaulting to: %v\n", address)
	}
	if os.Getenv("CALLBACKURL") != "" {
		callbackUrl = os.Getenv("CALLBACKURL")
	} else {
		callbackUrl = "127.0.0.1"
		log.Printf("No callback URL specified, defaulting to: %v\n", callbackUrl)
	}


	// Kafka
	if os.Getenv("KAFKAADDRESS") != "" {
		kafkaAddress = os.Getenv("KAFKAADDRESS")
	} else {
		log.Printf("No kafka server address specified\n")
		kafkaAddress = "127.0.0.1:9092"
	}

	// Twitch
	if os.Getenv("TWITCHCLIENTID") != "" {
		clientId = os.Getenv("TWITCHCLIENTID")
	} else {
		log.Printf("No twitch client ID specified\n")
	}
	if os.Getenv("TWITCHCLIENTSECRET") != "" {
		clientSecret = os.Getenv("TWITCHCLIENTSECRET")
	} else {
		log.Printf("No twitch client secret specified\n")
	}
}

func main() {
	writeNotification()
	log.Printf("Server starting")
	http.HandleFunc("/callback", callback)
	http.HandleFunc("/submit", submit)
	http.HandleFunc("/", home)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("Error %v\n", err)
	}
}
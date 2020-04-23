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
	addr string
	kafkaAddress string
)

type notificationMsg struct {
	link string
	body string
}

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

// Writes the webhook event notification to kafka broker
func writeNotification(msg string) error {
	err := kafkaWrite("twitch", msg)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Respond to twitch callback challenge for GET and notifications for POST
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
		} else {
			log.Printf("Callback Body: %s", body)
			// TODO asses serializing to json instead of string
			/*
				body := notificationMsg{`[<https://api.twitch.tv/helix/webhooks/hub>; rel=\"hub\", <https://api.twitch.tv/helix/users/follows?first=1&to_id=188951100>; rel=\"self\"]`,
					`{"data":[{"followed_at":"2020-04-22T22:26:18Z","from_id":"59480475","from_name":"eventualdecline","to_id":"188951100","to_name":"james65535"}]}`}
			*/
			str := fmt.Sprintf("%#v", body)
			writeErr := writeNotification(str)
			if writeErr != nil {
				log.Printf("Cannot do write notification: %s", writeErr)
			}
		}
	}
}

func home(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		log.Printf("Home GET request received!\n")

	} else {return}
}

func init() {
	// Webserver setup
	if os.Getenv("ADDRESS") != "" {
		addr = os.Getenv("ADDRESS")
		log.Printf("Webserver address: %v\n", addr)
	} else {
		addr = ":8080"
		log.Printf("No server address specified, defaulting to: %v\n", addr)
	}

	// Kafka setup
	if os.Getenv("KAFKAADDRESS") != "" {
		kafkaAddress = os.Getenv("KAFKAADDRESS")
		log.Printf("Using kafka server: %v\n", kafkaAddress)
	} else {
		kafkaAddress = "127.0.0.1:9092"
		log.Printf("No kafka server address specified, defaulting to: %v\n", kafkaAddress)
	}
}

func main() {
	log.Printf("Callback server starting\n")
	http.HandleFunc("/callback", callback)
	http.HandleFunc("/", home)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Error %v\n", err)
	}
}
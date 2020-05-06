package main

import (
	"github.com/james65535/twitchtest/pkg/twitch"
	"log"
	"os"
)

var (
	clientId string
	clientSecret string
	callbackUrl string
)

func init() {
	// Twitch setup
	if os.Getenv("TWITCHCLIENTID") != "" {
		clientId = os.Getenv("TWITCHCLIENTID")
	} else {
		clientId = "yssjfkvublum0gb5iel02puoyz7d3k"
		log.Printf("No twitch client ID specified\n")
	}
	if os.Getenv("TWITCHCLIENTSECRET") != "" {
		clientSecret = os.Getenv("TWITCHCLIENTSECRET")
	} else {
		clientSecret = "nhm00ahdhy3hrt07o0ivnxz8203cpj"
		log.Printf("No twitch client secret specified\n")
	}
	if os.Getenv("CALLBACKURL") != "" {
		callbackUrl = os.Getenv("CALLBACKURL")
	} else {
		callbackUrl = "http://35.189.39.205/callback"
		log.Printf("No twitch callback URL specified\n")
	}
}

func main() {
	// TODO move to env vars
	idUrl := "https://id.twitch.tv/oauth2"
	apiUrl := "https://api.twitch.tv/helix"

	// Create an Access Token
	s, _, err := twitch.GetAccessToken(
		clientId,
		clientSecret,
		idUrl+"/token",
		idUrl+"/validate")
	if err != nil {
		log.Fatalf("Access token error: %s", err)
	}
	log.Printf("Access token: %s\n", s)

	// Perform twitch API query
	// query(s,clientId, apiUrl + "/users?login=james65535"

	// Subscribe to a user follows update topic
	if err := twitch.SubscribeWebhook(
		s,
		clientId,
		apiUrl+"/webhooks/hub",
		apiUrl+"/users/follows?first=1&to_id=188951100",
		callbackUrl)
	err != nil {
		log.Printf("Subscribe error: %s\n", err)
	}

	// Check for current subscriptions
	if err := twitch.GetSubs(s,clientId,apiUrl + "/webhooks/subscriptions")
	err != nil {
		log.Printf("Subscribe error: %s\n", err)
	}
}

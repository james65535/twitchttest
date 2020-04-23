package main

import (
	"github.com/james65535/twitchtest/pkg/twitch"
	"log"
)

func main() {
	// TODO move to env vars
	clientId := "yssjfkvublum0gb5iel02puoyz7d3k"
	clientSecret := "nhm00ahdhy3hrt07o0ivnxz8203cpj"
	idUrl := "https://id.twitch.tv/oauth2"
	apiUrl := "https://api.twitch.tv/helix"
	callbackUrl := "http://35.189.39.205/callback"

	// Create an Access Token
	s, _, tokenErr := twitch.GetAccessToken(
		clientId,
		clientSecret,
		idUrl+"/token",
		idUrl+"/validate")
	if tokenErr != nil {
		log.Fatalf("Access token error: %s", tokenErr)
	} else {
		log.Printf("Access token: %s\n", s)
	}

	// Perform twitch API query
	// query(s,clientId, apiUrl + "/users?login=james65535"

	// Subscribe to a user follows update topic
	subErr := twitch.SubscribeWebhook(
		s,
		clientId,
		apiUrl+"/webhooks/hub",
		apiUrl+"/users/follows?first=1&to_id=188951100",
		callbackUrl)
	if subErr != nil {
		log.Printf("Subscribe error: %s\n", subErr)
	}

	// Check for current subscriptions
	subsErr := twitch.GetSubs(s,clientId,apiUrl + "/webhooks/subscriptions")
	if subsErr != nil {
		log.Printf("Subscribe error: %s\n", subsErr)
	}
}

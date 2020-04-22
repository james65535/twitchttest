package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"../../pkg/twitch"
)

func main() {
	clientId := "yssjfkvublum0gb5iel02puoyz7d3k"
	clientSecret := "nhm00ahdhy3hrt07o0ivnxz8203cpj"
	idUrl := "https://id.twitch.tv/oauth2"
	apiUrl := "https://api.twitch.tv/helix"

	s, _, _ := twitch.GetAccessToken(clientId, clientSecret, idUrl+"/token")
	fmt.Printf("access token: %s\n", s)
	// query(s,clientId, apiUrl + "/users?login=james65535"
	topicUrl := apiUrl + "/users/follows?first=1&to_id=188951100"
	subErr := twitch.SubscribeWebhook(s,clientId,apiUrl + "/webhooks/hub", topicUrl)
	if subErr != nil {
		log.Printf("Subscribe error: %s\n", subErr)
	}
/*
	subsErr := getSubs(s,clientId,apiUrl + "/webhooks/subscriptions")
	if subsErr != nil {
		log.Printf("Subscribe error: %s\n", subsErr)
	}
 */
}
package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type appAccessToken struct {
	Access_token string
	Expires_in   int
	Token_type   string
}

type webhookSub struct {
	Callback string `json:"hub.callback"`
	Mode string `json:"hub.mode"`
	Topic string `json:"hub.topic"`
	Lease int `json:"hub.lease_seconds"`
}

// TODO remove logging

// Obtains an access token
func GetAccessToken(id string, secret string, url string) (accesToken string, expires int, err error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", 0, err
	}

	values := req.URL.Query()
	values.Add("client_id", id)
	values.Add("client_secret", secret)
	values.Add("grant_type", "client_credentials")
	req.URL.RawQuery = values.Encode()

	// Send request
	client := &http.Client{}
	resp, postErr := client.Do(req)

	if postErr != nil {
		return "", 0, postErr
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return "", 0, readErr
	}

	if resp.StatusCode == 200 {
		var a appAccessToken
		parseErr := json.Unmarshal(body, &a)
		if parseErr != nil {
			return "", 0, parseErr
		}
		return a.Access_token, a.Expires_in, nil
	} else {
		return "", 0, fmt.Errorf(string(body))
	}
}

// Validates Access Tokens
func validate(accessToken string, url string) error {
	// Setup request
	req, requestErr := http.NewRequest("GET", url, nil)
	if requestErr != nil {
		return requestErr
	}
	req.Header.Add("Authorization", "OAuth " + accessToken)
	client := &http.Client{}

	// Send request
	resp, clientErr := client.Do(req)
	if clientErr != nil {
		return clientErr
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}

	// TODO handle error codes example: {"status":401,"message":"invalid access token"}
	if resp.StatusCode == 200 {
		return nil // {"client_id":"somestuff","scopes":[],"expires_in":5190956}
	} else {
		return fmt.Errorf(string(body))
	}
}

// Generic query
func Query(accessToken string, clientId string, url string) error {
	// Setup request
	req, requestErr := http.NewRequest("GET", url, nil)
	if requestErr != nil {
		return requestErr
	}
	req.Header.Add("Authorization", "Bearer " + accessToken)
	req.Header.Add("Client-ID", clientId)
	client := &http.Client{}

	// Send request
	resp, clientErr := client.Do(req)
	if clientErr != nil {
		return clientErr
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}
	log.Printf("query body: %s", body)

	// TODO handle error codes example: {"status":401,"message":"invalid access token"}
	if resp.StatusCode == 200 {
		return nil // {"client_id":"somestuff","scopes":[],"expires_in":5190956}
	} else {
		return fmt.Errorf(string(body))
	}
}

// Subscribe to a webhook
func SubscribeWebhook (accessToken string, clientId string, apiUrl string, topicUrl string, callbackUrl string) error {
	wh := webhookSub{
		callbackUrl,
		"subscribe",
		topicUrl,
		6400}  // lease was picked arbitrarily
	b, marshalErr := json.Marshal(&wh)
	if marshalErr != nil {
		return marshalErr
	}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer " + accessToken)
	req.Header.Add("Client-ID", clientId)
	req.Header.Add("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, postErr := client.Do(req)
	if postErr != nil {
		return postErr
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}
	if resp.StatusCode == 202 {
		log.Printf("Subscription body: %s", body)
		return nil
	} else {
		return fmt.Errorf(string(body))
	}
}

// Get webhook subscriptions
func GetSubs (accessToken string, clientId string, apiUrl string) error {
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer " + accessToken)
	req.Header.Add("Client-ID", clientId)

	// Send request
	client := &http.Client{}
	resp, getErr := client.Do(req)

	if getErr != nil {
		return getErr
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}
	if resp.StatusCode == 200 {
		log.Printf("Subscription body: %s", body)
		return nil
	} else {
		return fmt.Errorf(string(body))
	}
}

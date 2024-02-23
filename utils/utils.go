package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func SendSMS(to, message string) error {

	api_key := os.Getenv("SMS_API_KEY")

	// curl -X POST \
	// https://api.sandbox.africastalking.com/version1/messaging \
	// -H 'Accept: application/json' \
	// -H 'Content-Type: application/x-www-form-urlencoded' \
	// -H 'apiKey: MyAppApiKey' \
	// -d 'username=MyAppUsername&to=%2B254711XXXYYY,%2B254733YYYZZZ&message=Hello%20World!

	reqUrl := "https://api.sandbox.africastalking.com/version1/messaging"
	username := os.Getenv("SMS_USERNAME")

	data := url.Values{}
	data.Set("username", username)
	data.Set("to", to)
	data.Set("message", message)

	fmt.Printf("Sending with values: %v\n", data)

	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", api_key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(respBody))

	return nil
}

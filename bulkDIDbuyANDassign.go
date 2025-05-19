package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Check and load the .env file
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Fatalf("Error: .env file not found in the current directory")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Check if the required environment variables are set
	accountSid := checkEnvVariable("TWILIO_ACCOUNT_SID")
	authToken := checkEnvVariable("TWILIO_AUTH_TOKEN")
	trunkSid := checkEnvVariable("TWILIO_TRUNK_SID")
	bundleSid := checkEnvVariable("TWILIO_BUNDLE_SID")

	// Base URL for Twilio API
	baseURL := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/IncomingPhoneNumbers.json"

	for i := 0; i < 10; i++ {
		// Step 1: Search for an available UK number using HTTP GET request
		searchURL := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/AvailablePhoneNumbers/GB/Mobile.json"
		params := url.Values{}
		params.Set("PageSize", "1")

		client := &http.Client{}
		req, err := http.NewRequest("GET", searchURL+"?"+params.Encode(), nil)
		req.SetBasicAuth(accountSid, authToken)
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to search phone numbers: %v", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %v", err)
		}

		// Parse the response
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		// Get the first available phone number
		availableNumbers := result["available_phone_numbers"].([]interface{})
		if len(availableNumbers) == 0 {
			log.Fatalf("No available phone numbers found for request %d", i+1)
		}

		phoneNumber := availableNumbers[0].(map[string]interface{})["phone_number"].(string)

		// Step 2: Purchase the phone number and attach TrunkSid and BundleSid
		form := url.Values{}
		form.Add("PhoneNumber", phoneNumber)
		form.Add("TrunkSid", trunkSid)   // Attach TrunkSid here
		form.Add("BundleSid", bundleSid) // Attach BundleSid here

		req, err = http.NewRequest("POST", baseURL, bytes.NewBufferString(form.Encode()))
		req.SetBasicAuth(accountSid, authToken)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err = client.Do(req)
		if err != nil {
			log.Fatalf("Failed to purchase phone number: %v", err)
		}
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body during purchase: %v", err)
		}

		// Check for success in the response
		if resp.StatusCode != 201 {
			log.Fatalf("Failed to purchase phone number. Response: %s", body)
		}

		// Parse the purchase response to get the phone number SID
		var purchaseResult map[string]interface{}
		json.Unmarshal(body, &purchaseResult)

		phoneSID := purchaseResult["sid"].(string)

		// Final output: phone number, SID, and Trunk SID
		fmt.Printf("Purchased and assigned phone number %s with SID %s to trunk %s successfully.\n", phoneNumber, phoneSID, trunkSid)
	}
}

// Function to check if a specific environment variable is set
func checkEnvVariable(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Error: Environment variable %s is not set", key)
	}
	return value
}

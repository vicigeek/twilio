package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func main() {
	// Your Twilio account SID and Auth Token
	accountSid := "x"
	authToken := "x"

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	// Fetch all incoming phone numbers
	params := &openapi.ListIncomingPhoneNumberParams{}
	resp, err := client.Api.ListIncomingPhoneNumber(params)
	if err != nil {
		log.Fatalf("Failed to fetch incoming phone numbers: %v", err)
	}

	// Specify the CSV file name
	csvFile := "incoming_phone_numbers.csv"

	// Open the CSV file for writing
	file, err := os.Create(csvFile)
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the CSV header
	writer.Write([]string{"Phone Number", "Friendly Name", "Date Created", "Date Updated"})

	// Write the phone numbers to the CSV file
	for _, number := range resp {
		var phoneNumber, friendlyName, dateCreated, dateUpdated string

		if number.PhoneNumber != nil {
			phoneNumber = *number.PhoneNumber
		}
		if number.FriendlyName != nil {
			friendlyName = *number.FriendlyName
		}
		if number.DateCreated != nil {
			dateCreated = *number.DateCreated
		}
		if number.DateUpdated != nil {
			dateUpdated = *number.DateUpdated
		}

		writer.Write([]string{
			phoneNumber,
			friendlyName,
			dateCreated,
			dateUpdated,
		})
	}

	fmt.Printf("Incoming phone numbers have been written to %s\n", csvFile)
}


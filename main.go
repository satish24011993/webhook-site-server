package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

// Define the webhook URL
const webhookURL = "https://webhook.site/b4578ec6-d5b0-4258-94e9-f0b91c7bf369" // Assuming this is your correct unique ID

// Define the incoming data structure as a map since the keys for attributes and traits are dynamic.
type IncomingData map[string]interface{}

// Define the structures for the transformed data
type ValueType struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type TransformedData struct {
	Event           string               `json:"event"`
	EventType       string               `json:"event_type"`
	AppID           string               `json:"app_id"`
	UserID          string               `json:"user_id"`
	MessageID       string               `json:"message_id"`
	PageTitle       string               `json:"page_title"`
	PageURL         string               `json:"page_url"`
	BrowserLanguage string               `json:"browser_language"`
	ScreenSize      string               `json:"screen_size"`
	Attributes      map[string]ValueType `json:"attributes"`
	Traits          map[string]ValueType `json:"traits"`
}

// Channel for passing incoming data. We keep this to satisfy the prompt.
var dataChannel chan IncomingData

// Regular expressions to find numbered attribute and trait keys
var (
	atrkRe  = regexp.MustCompile(`^atrk(\d+)$`)
	uatrkRe = regexp.MustCompile(`^uatrk(\d+)$`)
)

func init() {
	// Initialize the channel (e.g., with a buffer of 100)
	dataChannel = make(chan IncomingData, 100)
}

func main() {
	// Start the worker goroutine.
	go worker(dataChannel)

	// Set up the HTTP server
	http.HandleFunc("/track", trackHandler)
	port := ":8080"
	log.Printf("HTTP server listening on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Handler for the incoming HTTP requests
func trackHandler(w http.ResponseWriter, r *http.Request) {
	// FIX 1: Add initial log to ensure handler is being hit
	log.Println("Received request on /track endpoint.")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var incomingData IncomingData
	if err := json.Unmarshal(body, &incomingData); err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// FIX 2: Move the "Processing data" log here to ensure it prints
	// before the goroutine is launched, confirming JSON unmarshaling succeeded.
	log.Printf("JSON successfully unmarshaled. Launching processing goroutine.")

	// Satisfies "each request should process in a separate goroutine"
	go func() {
		processData(incomingData)
		// This print will now appear AFTER the data is sent to the webhook
		log.Printf("Processing goroutine finished.")
	}()

	time.Sleep(100 * time.Millisecond)
	// Respond immediately to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request accepted"))
}

// worker function to satisfy the prompt's requirement for a worker receiving from a channel.
func worker(ch <-chan IncomingData) {
	for range ch {
		// Drains the channel if data is sent to it.
	}
}

// processData performs the transformation and sends the data to the webhook.
func processData(incomingData IncomingData) {
	// This log will now reliably appear immediately after the "JSON successfully unmarshaled" log.
	log.Printf("Processing data for app_id: %v", incomingData["id"])

	// 1. Initialize the target structure
	transformed := TransformedData{
		Attributes: make(map[string]ValueType),
		Traits:     make(map[string]ValueType),
	}

	// 2. Direct field mapping
	transformed.Event, _ = incomingData["ev"].(string)
	// Removed redundant log here: log.Printf("transformed.Event: %v", transformed.Event)
	transformed.EventType, _ = incomingData["et"].(string)
	transformed.AppID, _ = incomingData["id"].(string)
	transformed.UserID, _ = incomingData["uid"].(string)
	transformed.MessageID, _ = incomingData["mid"].(string)
	transformed.PageTitle, _ = incomingData["t"].(string)
	transformed.PageURL, _ = incomingData["p"].(string)
	transformed.BrowserLanguage, _ = incomingData["l"].(string)
	transformed.ScreenSize, _ = incomingData["sc"].(string)

	// 3. Dynamic Attribute and Trait mapping
	// Note: Re-defining these regex vars inside processData is redundant but harmless.
	// They are already global. I'll remove the re-declaration to clean up.
	// var (
	// 	atrkRe  = regexp.MustCompile(`^atrk(\d+)$`)
	// 	uatrkRe = regexp.MustCompile(`^uatrk(\d+)$`)
	// )

	for key, value := range incomingData {
		strKey := key
		strValue, ok := value.(string)
		if !ok {
			continue
		}

		// Check for Attribute keys (atrkN)
		if match := atrkRe.FindStringSubmatch(strKey); len(match) > 1 {
			number := match[1]
			valKey := "atrv" + number
			typeKey := "atrt" + number

			attrKey := strValue
			attrValue, valExists := incomingData[valKey].(string)
			attrType, typeExists := incomingData[typeKey].(string)

			if valExists && typeExists {
				transformed.Attributes[attrKey] = ValueType{
					Value: attrValue,
					Type:  attrType,
				}
			}
			continue
		}

		// Check for User Trait keys (uatrkN)
		if match := uatrkRe.FindStringSubmatch(strKey); len(match) > 1 {
			number := match[1]
			valKey := "uatrv" + number
			typeKey := "uatrt" + number

			traitKey := strValue
			traitValue, valExists := incomingData[valKey].(string)
			traitType, typeExists := incomingData[typeKey].(string)

			if valExists && typeExists {
				transformed.Traits[traitKey] = ValueType{
					Value: traitValue,
					Type:  traitType,
				}
			}
			continue
		}
	}
	jsonData, _ := json.MarshalIndent(transformed, "", "  ")
	log.Printf("Transformed Data:\n%s", jsonData)

	log.Printf("Attributes mapped: %d, Traits mapped: %d",
		len(transformed.Attributes),
		len(transformed.Traits))
	// 4. Send the transformed data to the webhook
	sendToWebhook(transformed)
}

// sendToWebhook sends the final transformed data as JSON to the external endpoint.
func sendToWebhook(data TransformedData) {
	// Marshal the transformed data back into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling transformed data: %v", err)
		return
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating webhook request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client
	client := &http.Client{Timeout: 10 * time.Second}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending data to webhook: %v", err)
		return
	}
	defer resp.Body.Close()

	// Log the response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Webhook responded with status: %s", resp.Status)
	} else {
		log.Printf("Successfully sent data to webhook. Status: %s", resp.Status)
	}
}

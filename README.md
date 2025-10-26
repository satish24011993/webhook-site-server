🚀 Go Webhook Data Transformer

This project implements a single-file Go HTTP server designed to process incoming webhook payloads. It converts a flat JSON structure with abbreviated, numbered keys (e.g., ev, atrk1) into a rich, nested JSON format before forwarding the result to an external webhook endpoint.

✨ Features

Endpoint: Handles incoming POST requests at /track.

Asynchronous Processing: Uses Go goroutines for non-blocking I/O.

Data Transformation: Maps short-key data fields to full, descriptive fields and structures dynamic attributes (atrkN) and traits (uatrkN) into nested maps.

Reliable Delivery: Forwards transformed data using an HTTP client with a 10-second timeout.

🛠️ Getting Started

Prerequisites

Go: Go 1.18+ installed on your system.

cURL or Postman: For sending test requests.

Setup and Configuration

Save the Code: Ensure the provided Go code is saved as a single file named main.go.

Configure Webhook URL:

Before running, you MUST update the placeholder webhookURL constant in main.go.

Go to https://webhook.site/ to generate a new, unique URL, and replace the value in main.go:

// main.go (Line 15)
const webhookURL = "[https://webhook.site/YOUR-NEW-UNIQUE-ID-GOES-HERE](https://webhook.site/YOUR-NEW-UNIQUE-ID-GOES-HERE)" 


Run the Server:
Open your terminal in the project directory and execute:

go run main.go


The console will display: HTTP server listening on port :8080.

🧪 Testing the Transformation

While the server is running, use cURL or Postman in a separate window to send the test payload to your local server.

1. Input Sample

This JSON will be sent to your Go server:

{
  "ev": "contact_form_submitted",
  "et": "form_submit",
  "id": "cl_app_id_001",
  "uid": "cl_app_id_001-uid-001",
  "mid": "cl_app_id_001-uid-001",
  "t": "Vegefoods - Free Bootstrap 4 Template by Colorlib",
  "p": "[http://shielded-eyrie-45679.herokuapp.com/contact-us](http://shielded-eyrie-45679.herokuapp.com/contact-us)",
  "l": "en-US",
  "sc": "1920 x 1080",
  "atrk1": "form_varient",
  "atrv1": "red_top",
  "atrt1": "string",
  "atrk2": "ref",
  "atrv2": "XPOWJRICW993LKJD",
  "atrt2": "string",
  "uatrk1": "name",
  "uatrv1": "iron man",
  "uatrt1": "string",
  "uatrk2": "email",
  "uatrv2": "ironman@avengers.com",
  "uatrt2": "string",
  "uatrk3": "age",
  "uatrv3": "32",
  "uatrt3": "integer"
}


2. cURL Command

Use this command to send the request (ensure the JSON is correctly escaped for the shell):

curl -X POST http://localhost:8080/track \
-H "Content-Type: application/json" \
-d '{"ev":"contact_form_submitted","et":"form_submit","id":"cl_app_id_001","uid":"cl_app_id_001-uid-001","mid":"cl_app_id_001-uid-001","t":"Vegefoods - Free Bootstrap 4 Template by Colorlib","p":"[http://shielded-eyrie-45679.herokuapp.com/contact-us](http://shielded-eyrie-45679.herokuapp.com/contact-us)","l":"en-US","sc":"1920 x 1080","atrk1":"form_varient","atrv1":"red_top","atrt1":"string","atrk2":"ref","atrv2":"XPOWJRICW993LKJD","atrt2":"string","uatrk1":"name","uatrv1":"iron man","uatrt1":"string","uatrk2":"email","uatrv2":"ironman@avengers.com","uatrt2":"string","uatrk3":"age","uatrv3":"32","uatrt3":"integer"}'


3. Verification

Go Terminal: Check your terminal for successful logs, including the full transformed JSON printed by the temporary debug log:

Transformed Data:
{... (full nested JSON structure) ...}
Successfully sent data to webhook. Status: 200 OK


Webhook Site: Visit your unique webhook URL. Click on the latest POST request in the left sidebar and inspect the Request Body to view the transformed JSON.

📊 Expected Output Structure

The data sent to the external webhook will be in this final, required format:

{
  "event": "contact_form_submitted",
  "event_type": "form_submit",
  "app_id": "cl_app_id_001",
  "user_id": "cl_app_id_001-uid-001",
  "message_id": "cl_app_id_001-uid-001",
  "page_title": "Vegefoods - Free Bootstrap 4 Template by Colorlib",
  "page_url": "[http://shielded-eyrie-45679.herokuapp.com/contact-us](http://shielded-eyrie-45679.herokuapp.com/contact-us)",
  "browser_language": "en-US",
  "screen_size": "1920 x 1080",
  "attributes": {
    "form_varient": {
      "value": "red_top",
      "type": "string"
    },
    "ref": {
      "value": "XPOWJRICW993LKJD",
      "type": "string"
    }
  },
  "traits": {
    "name": {
      "value": "iron man",
      "type": "string"
    },
    "email": {
      "value": "ironman@avengers.com",
      "type": "string"
    },
    "age": {
      "value": "32",
      "type": "integer"
    }
  }
}

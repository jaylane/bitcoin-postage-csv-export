package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type OrderList struct {
	OrderID        string `json:"order_id"`
	OrderTimestamp string `json:"order_timestamp"`
	Price          string `json:"price"`
}

type Label struct {
	From       string `json:"from"`
	To         string `json:"to"`
	TrackingNo string `json:"tracking_no"`
	ShipmentID string `json:"shipment_id"`
	Carrier    string `json:"carrier"`
}

func makePostRequest(url string, apiKey string, apiSecret string, orderID string) ([]byte, error) {
	client := &http.Client{}

	// Create form data
	formData := fmt.Sprintf("key=%s&secret=%s", apiKey, apiSecret)
	if orderID != "" {
		formData += fmt.Sprintf("&order_id=%s", orderID)
	}

	// Create request with form data
	req, err := http.NewRequest("POST", url, strings.NewReader(formData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set content type for form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func main() {
	// Command line flags for API credentials
	apiKey := flag.String("key", "", "Bitcoin Postage API Key")
	apiSecret := flag.String("secret", "", "Bitcoin Postage API Secret")
	outputFile := flag.String("output", "orders.csv", "Output CSV file name")
	flag.Parse()

	// Validate API credentials
	if *apiKey == "" || *apiSecret == "" {
		fmt.Println("Error: API key and secret are required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Get list of orders
	body, err := makePostRequest("https://bitcoinpostage.info/api/orders", *apiKey, *apiSecret, "")
	if err != nil {
		fmt.Printf("Error getting orders list: %v\n", err)
		os.Exit(1)
	}

	// Parse JSON response for order list
	var orders []OrderList
	if err := json.Unmarshal(body, &orders); err != nil {
		fmt.Printf("Error parsing orders JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d orders\n", len(orders))

	// Create CSV file
	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	headers := []string{
		"Order ID",
		"Order Timestamp",
		"Date Time",
		"Total Price",
		"From",
		"To",
		"Tracking Number",
		"Shipment ID",
		"Carrier",
	}
	if err := writer.Write(headers); err != nil {
		fmt.Printf("Error writing CSV headers: %v\n", err)
		os.Exit(1)
	}

	// Process each order
	for _, order := range orders {
		fmt.Printf("Fetching details for order %s...\n", order.OrderID)

		// Get order details
		body, err := makePostRequest("https://bitcoinpostage.info/api/retrieve-order", *apiKey, *apiSecret, order.OrderID)
		if err != nil {
			fmt.Printf("Error getting order details for %s: %v\n", order.OrderID, err)
			continue
		}

		var labels []Label
		if err := json.Unmarshal(body, &labels); err != nil {
			fmt.Printf("Error parsing order detail JSON for %s: %v\n", order.OrderID, err)
			continue
		}

		// Convert Unix timestamp to readable date
		timestamp, err := strconv.ParseInt(order.OrderTimestamp, 10, 64)
		if err != nil {
			fmt.Printf("Error parsing timestamp for order %s: %v\n", order.OrderID, err)
			continue
		}
		dateTime := time.Unix(timestamp, 0).Format(time.RFC3339)

		// Write each label to CSV
		for _, label := range labels {
			record := []string{
				order.OrderID,
				order.OrderTimestamp,
				dateTime,
				order.Price,
				label.From,
				label.To,
				label.TrackingNo,
				label.ShipmentID,
				label.Carrier,
			}
			if err := writer.Write(record); err != nil {
				fmt.Printf("Error writing order to CSV: %v\n", err)
				continue
			}
		}
	}

	fmt.Printf("Successfully exported orders to %s\n", *outputFile)
}

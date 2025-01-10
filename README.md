# Bitcoin Postage Order Export Tool

This tool fetches all historical orders from your Bitcoin Postage account and exports them to a CSV file.

## Usage

```bash
go run fetch_orders.go -key YOUR_API_KEY -secret YOUR_API_SECRET [-output orders.csv]
```

### Parameters:

- `-key`: Your Bitcoin Postage API Key (required)
- `-secret`: Your Bitcoin Postage API Secret (required)
- `-output`: Output CSV filename (optional, defaults to "orders.csv")

### Output Format

The script will create a CSV file with the following columns:
- Order ID
- Timestamp (Unix timestamp)
- Date Time (Human readable format)
- Price

## Requirements

- Go 1.x
- Bitcoin Postage API credentials

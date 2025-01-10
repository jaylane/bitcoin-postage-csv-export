# Bitcoin Postage Order Export Tool

This tool fetches all historical orders from your Bitcoin Postage account and exports them to a CSV file. For each order, it retrieves detailed shipping information including tracking numbers and addresses.

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
- Order Timestamp (Unix timestamp)
- Date Time (Human readable format in RFC3339)
- Total Price
- From Address
- To Address
- Tracking Number
- Shipment ID
- Carrier

## Requirements

- Go 1.16 or later
- Bitcoin Postage API credentials

## API Endpoints Used

The tool interacts with the following Bitcoin Postage API endpoints:
- `https://bitcoinpostage.info/api/orders` - Retrieves list of all orders
- `https://bitcoinpostage.info/api/retrieve-order` - Retrieves detailed information for a specific order

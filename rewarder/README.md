# Rewarder App

## API Endpoints

### Upload CSV File
- **POST** `/api/upload`
- **Content-Type**: `multipart/form-data`
- **Body**: CSV file with the field name `file`
- **Response**: JSON with success status and generated vouchers

## CSV Format

The application expects CSV files with the following columns (in order):

| Column | Description | Example |
|--------|-------------|---------|
| Customer ID | Unique customer identifier | `1` |
| Customer First Name | Customer's first name | `John` |
| Order Value | Total order value | `1500` |

### Example CSV Content
```csv
Customer ID,Customer First Name,Order Value
1,Kweku,100
2,Abena,1200
3,Kojo,4800
4,Esi,7500
5,Yaw,15000
6,Akua,999
7,Mensah,10000
```

## Installation & Setup

### Prerequisites
- Go 1.24.4 or later
- Git

### Installation

1. **Clone the repository**
   ```bash
    git clone <repository-url>
    cd rewarder
   ```

2. **Install dependencies**
   ```bash
    go mod tidy
   ```

3. **Run the application**
   ```bash
    go run .
   ```

## Usage Examples

### Upload a CSV File

```bash
  curl -X POST http://localhost:8080/api/upload \
    -F "file=@orders.csv"
```

**Response:**
```json
{
  "message": "file uploaded successfully",
  "data": [
    {
      "ID": "4096d289-6178-4ff1-ae71-f24413c10b39",
      "CustomerID": 2,
      "CustomerName": "abena",
      "OrderValue": 1200,
      "Amount": 100,
      "CreatedAt": "2025-09-19T22:55:42.112955Z",
      "ExpiresAt": "2025-09-20T22:55:42.112955Z"
    },
    {
      "ID": "168f9188-2f30-4438-a611-d2dd1bb35f21",
      "CustomerID": 3,
      "CustomerName": "kojo","OrderValue": 4800,
      "Amount": 100,
      "CreatedAt": "2025-09-19T22:55:42.11626Z",
      "ExpiresAt": "2025-09-20T22:55:42.11626Z"
    },
  ],
  "success": true
}
```

## Voucher Structure

Each generated voucher contains:

- **ID**: Unique UUID identifier
- **CustomerID**: Original customer ID from CSV
- **CustomerName**: Customer's first name (lowercased)
- **OrderValue**: Original order value that triggered the voucher
- **Amount**: Voucher value based on reward tier
- **CreatedAt**: Timestamp when voucher was created
- **ExpiresAt**: Timestamp when voucher expires

## Testing

Run the test suite:

```bash
  go test -v ./...
```


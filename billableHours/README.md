# Billable Hours App

## API Endpoints

### Upload CSV File
- **POST** `/api/upload`
- **Content-Type**: `multipart/form-data`
- **Body**: CSV file with the field name `file`
- **Response**: JSON with success status and processed data

### Download Invoice
- **GET** `/api/download/{companyName}`
- **Response**: PDF file download
- **Content-Type**: `application/pdf`

## CSV Format

The application expects CSV files with the following columns (in order):

| Column | Description | Example |
|--------|-------------|---------|
| Employee ID | Unique employee identifier | `1` |
| Billable Rate (per hour) | Hourly rate in currency | `100.00` |
| Project | Company/Project name | `Acme Corp` |
| Date | Work date | `2019-07-01` |
| Start time | Work start time (24-hour format) | `09:00` |
| End Time | Work end time (24-hour format) | `17:00` |

### Example CSV Content
```csv
"Employee ID","Billable Rate (per hour)","Project","Date","Start time","End Time"
"1","100","Acme","2019-07-01","09:00","11:00"
"1","100","Acme","2019-07-01","12:00","14:00"
"2","150","Acme","2019-07-01","10:00","15:00"
```

## Installation & Setup

### Prerequisites
- Go 1.24.4 or later
- Git

### Installation

1. **Clone the repository**
   ```bash
    git clone <repository-url>
    cd billableHours
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
    -F "file=@timesheet.csv"
```

**Response:**
```json
{
  "message": "file uploaded successfully",
  "data": {
    "acme": {
      "1": {
        "BillableRate": 100,
        "TotalHours": 4
      },
      "2": {
        "BillableRate": 150,
        "TotalHours": 5
      }
    }
  },
  "success": true
}
```

### Download an Invoice

```bash
  curl -X GET http://localhost:8080/api/download/Acme \
    --output invoice.pdf
```


## Testing

Run the test suite:

```bash
  go test -v ./...
```


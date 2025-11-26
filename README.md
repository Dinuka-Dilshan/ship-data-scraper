# Ship Data Scraper

A Go-based serverless function that scrapes vessel and schedule data from the Maersk API, processes it in parallel, and exports results as CSV files in a ZIP archive.

## Overview

This project is built as a Google Cloud Function that:

- Fetches active vessels from the Maersk API based on carrier codes
- Retrieves vessel schedules for each ship within a specified date range
- Processes schedule data concurrently for improved performance
- Generates two CSV files:
  - `data.csv` - Complete vessel schedule information
  - `nodes.csv` - Port/node frequency analysis
- Returns results as a downloadable ZIP file

## Features

- **Concurrent Processing**: Uses goroutines and WaitGroups to fetch and process vessel schedules in parallel
- **CSV Export**: Structured data export with comprehensive vessel and schedule information
- **HTTP-Based API**: Query parameter configuration for flexible usage
- **Error Handling**: Comprehensive logging with structured logging (slog)
- **ZIP Distribution**: Multiple CSV files packaged for easy distribution

## Project Structure

```
.
├── main.go                 # HTTP handler and orchestration logic
├── maersk.go              # Maersk API client and data models
├── app_http_client.go     # HTTP client wrapper for API calls
├── csv.go                 # CSV formatting and export utilities
├── utils.go               # Helper functions
├── go.mod                 # Go module dependencies
└── README.md              # This file
```

## API Endpoint

**Query Parameters**:

- `apiKey` (required) - Maersk API authentication key
- `carrierCodes` (required) - Comma-separated carrier codes
- `startDate` (required) - Schedule start date (format: YYYY-MM-DD)
- `endDate` (required) - Schedule end date (format: YYYY-MM-DD)

**Example Request**:

```
GET https://function-url.cloudfunctions.net/Handler?apiKey=YOUR_API_KEY&carrierCodes=MSC,MAEU&startDate=2025-01-01&endDate=2025-12-31
```

**Response**:

- Content-Type: `application/zip`
- Contains `data.csv` and `nodes.csv`

## Output Files

### data.csv

Contains complete vessel schedule information with columns:

- Vessel Name, Call Sign, Flag ISO Country Code, IMO Number, Maersk Code
- Port: ISO Country Code, Country Name, UN Location Code, City, Port Name, Port Code, Region Code
- Terminal: Name, RKST Code, Geo Code
- Arrival: Time, Timing Classifier, Voyage Number, Service Name, Service Code
- Departure: Time, Timing Classifier, Voyage Number, Service Name, Service Code

### nodes.csv

Contains port/node frequency analysis:

- Node Name
- Frequency (number of scheduled visits)




# GetUsage Sample

This sample demonstrates how to retrieve usage data from Temporal Cloud using the Cloud API.

## Prerequisites

- Temporal Cloud account
- Cloud API key with appropriate permissions

## Setup

1. Set your Temporal Cloud API key as an environment variable:
   ```bash
   export TEMPORAL_CLOUD_API_KEY="your-api-key-here"
   ```

## Running the Sample

```bash
go run main.go
```

## What it does

The sample:
1. Creates a Temporal Cloud client using your API key
2. Requests usage data for the last 30 days
3. Displays summary usage information
4. Shows usage breakdown by namespace if available

## Output

The sample will display:
- Usage summaries for the specified time range
- Each summary contains detailed usage information
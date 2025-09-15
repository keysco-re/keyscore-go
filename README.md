# keyscore-go

A Go client library for the Keysco.re API, providing access to comprehensive threat intelligence data including hash lookups, IP geolocation, machine information, and advanced search capabilities.

## Features

- **Hash Lookup**: Query hash databases for malware analysis
- **IP Lookup**: Get geolocation and ISP information for IP addresses
- **Search**: Advanced search across multiple data sources
- **Count**: Get result counts for queries
- **Machine Info**: Retrieve detailed machine information by UUID
- **Download**: Download files and archives by UUID
- **Sources**: List available data sources
- **Health Check**: Monitor API service status

## Installation

```bash
go get github.com/keysco-re/keyscore-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/keysco-re/keyscore-go/keyscore"
)

func main() {
    // Create a new client with API key
    client := keyscore.NewClient(
        keyscore.WithAPIKey("your-api-key-here"),
    )
    
    ctx := context.Background()
    
    // Perform a hash lookup
    hashResp, err := client.HashLookup(ctx, keyscore.HashLookupRequest{
        Terms: []string{"5d41402abc4b2a76b9719d911017c592"},
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d results\n", hashResp.Size)
    for hash, record := range hashResp.Results {
        fmt.Printf("Hash: %s, Type: %s, Source: %s\n", 
            hash, record.Type, record.Source)
    }
}
```

## Authentication

The client supports API key authentication:

```go
// Using API key
client := keyscore.NewClient(
    keyscore.WithAPIKey("your-api-key-here"),
)

// Custom base URL (for enterprise installations)
client := keyscore.NewClient(
    keyscore.WithAPIKey("your-api-key-here"),
    keyscore.WithBaseURL("https://your-custom-endpoint.com"),
)

// Custom HTTP client with timeout
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}
client := keyscore.NewClient(
    keyscore.WithAPIKey("your-api-key-here"),
    keyscore.WithHTTPClient(httpClient),
)
```

## API Reference

### Health Check

Check the API service health:

```go
health, err := client.Health(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("API Status: %s\n", health.Status)
```

### Sources

List available data sources:

```go
sources, err := client.Sources(ctx)
if err != nil {
    log.Fatal(err)
}

for name, info := range sources.Sources {
    fmt.Printf("Source: %s, Display Name: %s\n", name, info.DisplayName)
    fmt.Printf("Allowed Types: %v\n", info.AllowedTypes)
}
```

### Hash Lookup

Lookup hash information:

```go
resp, err := client.HashLookup(ctx, keyscore.HashLookupRequest{
    Terms: []string{
        "5d41402abc4b2a76b9719d911017c592",
        "098f6bcd4621d373cade4e832627b4f6",
    },
})
if err != nil {
    log.Fatal(err)
}

for hash, record := range resp.Results {
    fmt.Printf("Hash: %s\n", hash)
    fmt.Printf("  Type: %s\n", record.Type)
    fmt.Printf("  Plaintext: %s\n", record.Plaintext)
    fmt.Printf("  Source: %s\n", record.Source)
    fmt.Printf("  First Seen: %s\n", record.FirstSeen)
}
```

### IP Lookup

Get IP geolocation information:

```go
resp, err := client.IPLookup(ctx, keyscore.IPLookupRequest{
    Terms: []string{"8.8.8.8", "1.1.1.1"},
})
if err != nil {
    log.Fatal(err)
}

for ip, info := range resp.Results {
    fmt.Printf("IP: %s\n", ip)
    fmt.Printf("  Country: %s (%s)\n", info.Country, info.CountryCode)
    fmt.Printf("  City: %s\n", info.City)
    fmt.Printf("  ISP: %s\n", info.ISP)
    fmt.Printf("  Coordinates: %.4f, %.4f\n", info.Lat, info.Lon)
}

// Check for errors
for ip, errMsg := range resp.Errors {
    fmt.Printf("Error for IP %s: %s\n", ip, errMsg)
}
```

### Search

Perform advanced searches:

```go
resp, err := client.Search(ctx, keyscore.SearchRequest{
    Terms:    []string{"malware", "trojan"},
    Types:    []string{"hash", "domain"},
    Source:   "virustotal",
    Wildcard: true,
    Page:     1,
    PageSize: 100,
    DateFrom: "2023-01-01",
    DateTo:   "2023-12-31",
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Search took %dms, found %d results\n", resp.Took, resp.Size)
for source, results := range resp.Results {
    fmt.Printf("Source %s: %d results\n", source, len(results))
}
```

### Count

Get result counts without retrieving full data:

```go
// Basic count
count, err := client.Count(ctx, keyscore.CountRequest{
    Terms: []string{"malware"},
    Types: []string{"hash"},
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Total count: %d\n", count.Count)

// Detailed count with breakdown
detailed, err := client.CountDetailed(ctx, keyscore.CountRequest{
    Terms:    []string{"malware"},
    Types:    []string{"hash", "domain"},
    Wildcard: true,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total: %d (took %dms)\n", detailed.TotalCount, detailed.Took)
for source, count := range detailed.Counts {
    fmt.Printf("  %s: %d\n", source, count)
}
```

### Machine Information

Retrieve machine information by UUID:

```go
machine, err := client.MachineInfo(ctx, "550e8400-e29b-41d4-a716-446655440000")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Machine Info:\n")
fmt.Printf("  Computer Name: %s\n", machine.ComputerName)
fmt.Printf("  User: %s\n", machine.UserName)
fmt.Printf("  OS: %s\n", machine.OperationSystem)
fmt.Printf("  IP: %s\n", machine.IP)
fmt.Printf("  Country: %s\n", machine.Country)
fmt.Printf("  RAM: %s\n", machine.InstalledRAM)
fmt.Printf("  Processor: %s\n", machine.Processor)
```

### Download Files

Download files by UUID:

```go
// Download full archive
result, err := client.Download(ctx, "550e8400-e29b-41d4-a716-446655440000", "")
if err != nil {
    log.Fatal(err)
}
defer result.Body.Close()

fmt.Printf("Content Type: %s\n", result.ContentType)
fmt.Printf("Content Length: %d bytes\n", result.ContentLength)

// Save to file
file, err := os.Create("download.zip")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

_, err = io.Copy(file, result.Body)
if err != nil {
    log.Fatal(err)
}

// Download specific file
fileResult, err := client.Download(ctx, "550e8400-e29b-41d4-a716-446655440000", "path/to/file.txt")
if err != nil {
    log.Fatal(err)
}
defer fileResult.Body.Close()

// Read file content
content, err := io.ReadAll(fileResult.Body)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("File content: %s\n", string(content))
```

## Data Models

### Core Types

- `HashRecord`: Hash lookup result with plaintext, type, source, and timestamp
- `IPInfo`: IP geolocation data including country, city, ISP, and coordinates
- `MachineInfo`: Detailed machine information including hardware and system details
- `SourceInfo`: Data source metadata with allowed types and capabilities

### Request Types

- `HashLookupRequest`: Hash lookup parameters
- `IPLookupRequest`: IP lookup parameters
- `SearchRequest`: Advanced search parameters with filtering options
- `CountRequest`: Count query parameters

### Response Types

- `HashLookupResponse`: Hash lookup results with timing information
- `IPLookupResponse`: IP lookup results with error handling
- `SearchResponse`: Search results with pagination
- `CountResponse`: Simple count result
- `DetailedCountResponse`: Count breakdown by source

## Error Handling

The client returns structured errors for API failures:

```go
resp, err := client.HashLookup(ctx, keyscore.HashLookupRequest{
    Terms: []string{"invalid-hash"},
})
if err != nil {
    if apiErr, ok := err.(*keyscore.APIError); ok {
        fmt.Printf("API Error: %d - %s\n", apiErr.StatusCode, apiErr.Message)
        switch apiErr.StatusCode {
        case 401:
            fmt.Println("Authentication failed - check your API key")
        case 429:
            fmt.Println("Rate limit exceeded - please wait before retrying")
        case 500:
            fmt.Println("Server error - please try again later")
        }
    } else {
        fmt.Printf("Network error: %v\n", err)
    }
    return
}
```

## Examples

### Batch Hash Analysis

```go
func analyzeHashes(client *keyscore.Client, hashes []string) {
    ctx := context.Background()
    
    // Split into batches of 100
    batchSize := 100
    for i := 0; i < len(hashes); i += batchSize {
        end := i + batchSize
        if end > len(hashes) {
            end = len(hashes)
        }
        
        batch := hashes[i:end]
        resp, err := client.HashLookup(ctx, keyscore.HashLookupRequest{
            Terms: batch,
        })
        if err != nil {
            log.Printf("Batch %d failed: %v", i/batchSize+1, err)
            continue
        }
        
        fmt.Printf("Batch %d: %d/%d hashes found\n", 
            i/batchSize+1, resp.Size, len(batch))
    }
}
```

### Threat Intelligence Pipeline

```go
func threatIntelligence(client *keyscore.Client, indicators []string) {
    ctx := context.Background()
    
    // First, get counts to understand data volume
    count, err := client.CountDetailed(ctx, keyscore.CountRequest{
        Terms: indicators,
        Types: []string{"hash", "domain", "ip"},
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Total indicators found: %d\n", count.TotalCount)
    
    // If manageable volume, get full results
    if count.TotalCount < 10000 {
        search, err := client.Search(ctx, keyscore.SearchRequest{
            Terms:    indicators,
            Types:    []string{"hash", "domain", "ip"},
            PageSize: 1000,
        })
        if err != nil {
            log.Fatal(err)
        }
        
        // Process results
        for source, results := range search.Results {
            fmt.Printf("Processing %d results from %s\n", len(results), source)
            // Your analysis logic here
        }
    }
}
```

## Requirements

- Go 1.19 or higher
- Valid Keysco.re API key
- Internet connection for API access

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions:

- **Documentation**: [https://docs.keysco.re](https://docs.keysco.re)
- **Issues**: [GitHub Issues](https://github.com/keysco-re/keyscore-go/issues)
- **Email**: esson@riseup.net

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
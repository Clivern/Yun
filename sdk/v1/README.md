# Mut SDK v1

A Go SDK for interacting with the Mut API.

## Installation

```bash
go get github.com/clivern/mut/sdk/v1
```

## Usage

### Basic Client Setup

```go
package main

import (
    "fmt"
    "log"

    "github.com/clivern/mut/sdk/v1"
)

func main() {
    // Create a new client
    client, err := v1.NewClient(v1.ClientConfig{
        BaseURL: "https://api.example.com",
        APIKey:  "your-api-key", // Optional
    })
    if err != nil {
        log.Fatal(err)
    }

    // Use the client
    health, err := client.Health()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Health: %s\n", health.Status)
}
```

### Standalone Functions

You can also use standalone functions without creating a client:

```go
// Health check
health, err := v1.Health("https://api.example.com")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Health: %s\n", health.Status)

// Readiness check
ready, err := v1.Ready("https://api.example.com")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Ready: %s\n", ready.Status)
```

## Available Endpoints

### Health Check

The health check endpoint verifies that the API is running.

```go
health, err := client.Health()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Status: %s\n", health.Status) // Status: ok
```

### Readiness Check

The readiness check endpoint verifies that the API is ready to serve traffic (checks database connectivity).

```go
ready, err := client.Ready()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Status: %s\n", ready.Status) // Status: ok
```

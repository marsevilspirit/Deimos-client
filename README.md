# Deimos Client

<p align="center">
  <strong>The Go client for Deimos</strong>
</p>

---

## Introduction

`Deimos-client` is the Go client for the [Deimos](https://github.com/marsevilspirit/deimos) distributed key-value store. It provides a simple and easy-to-use API that allows you to easily interact with a Deimos cluster in your own Go applications.

## Features

- **Concise API**: Provides an intuitive, chainable API for key-value operations.
- **Automatic Node Discovery**: Automatically manages nodes in the Deimos cluster and is aware of node changes.
- **Load Balancing**: Automatically load balances requests across multiple Deimos nodes.
- **Type-Safe**: Provides a type-safe interface to reduce runtime errors.
- **Extensible**: Easy to extend, with support for custom HTTP clients and other configurations.

## Installation

```bash
go get github.com/marsevilspirit/deimos-client
```

## Quick Start

Here is a simple example of how to use `Deimos-client` to connect to a Deimos cluster and perform basic operations.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	deimosclient "github.com/marsevilspirit/deimos-client"
)

func main() {
	// Create a new Deimos client
	client := deimosclient.New(deimosclient.WithEndpoints("http://127.0.0.1:4001"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set a key-value pair
	setResp, err := client.Set(ctx, "/mykey", "hello world")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Set Response: %+v\n", setResp)

	// Get the value of a key
	getResp, err := client.Get(ctx, "/mykey")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Get Response: %+v\n", getResp)

	// Delete a key
	delResp, err := client.Delete(ctx, "/mykey")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Delete Response: %+v\n", delResp)
}
```

## API Usage Examples

### Create Client

```go
// Create a client connected to a single node
client := deimosclient.New(deimosclient.WithEndpoints("http://127.0.0.1:4001"))

// Create a client connected to multiple nodes
client = deimosclient.New(deimosclient.WithEndpoints("http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"))
```

### Set Key-Value

```go
// Set a simple key-value pair
resp, err := client.Set(ctx, "/foo", "bar")

// Set a key-value pair with a TTL (Time-To-Live)
resp, err = client.Set(ctx, "/foo", "bar", deimosclient.WithTTL(10*time.Second))
```

### Get Key-Value

```go
// Get the value of a single key
resp, err := client.Get(ctx, "/foo")

// Recursively get all key-values under a directory
resp, err = client.Get(ctx, "/dir", deimosclient.WithRecursive())
```

### Delete Key-Value

```go
// Delete a single key
resp, err := client.Delete(ctx, "/foo")

// Recursively delete a directory
resp, err = client.Delete(ctx, "/dir", deimosclient.WithRecursive())
```

### Watch for Key Changes

```go
// Watch for changes on a single key
watcher := client.Watch(ctx, "/foo")
for resp := range watcher {
    fmt.Printf("Key changed: %+v\n", resp)
}

// Recursively watch for changes in a directory
watcher = client.Watch(ctx, "/dir", deimosclient.WithRecursive())
for resp := range watcher {
    fmt.Printf("Key in directory changed: %+v\n", resp)
}
```

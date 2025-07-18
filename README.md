# Deimos Client: Your Go Gateway to the Deimos Ecosystem

<div align="center">
  <pre>
⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⠇⠀⠘⢿⣿⣿⣧⣶⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣤⣤⣾⣿⣿⣿⣿⣿⣿⣿⣿⡿⣿⣿⣿⣿⣿⣿⣿⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠉⢿⣿⣿⣿⠿⠟⠋⠉⠀⠀⢸⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠸⣿⣿⣷⣶⣤⣤⣀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣧⣤⣤⣤⡀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠈⠙⠛⠛⠛⢿⣿⣧⠀⢸⣿⣿⣿⣿⣿⣿⣿⠟⠛⠻⣿⣦⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⣴⣶⣶⡇⠀⠀⠀⣿⣿⠀⢼⣿⣿⣿⣿⣿⣿⣿⣤⡀⠀⢻⣿⡇⠀⠀⠀⠀⠀
⠀⠀⠀⠈⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⠀⠉⠉⠉⠀⠀⠀⢸⣿⡟⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⠃⣶⣾⣿⡇⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠈⠉⠙⢻⣿⣿⣿⣿⣿⣿⣿⠀⠀⢸⣿⣿⠇⠀⢸⣿⡇⠀⢸⣿⣿⣿⣿⣿⣿⣿⡏⠀⠛⣿⣿⠇⠀⠀⠀⠀⠀
⣤⣤⣤⣤⣤⣤⣤⣤⣼⣿⣿⣿⣿⣿⣿⣿⣷⣤⣤⣤⣤⣤⣤⣼⣿⣧⣤⣤⣿⣿⣿⣿⣿⣿⣿⣧⣤⣴⣿⣿⣤⣤⣤⣤⣤⣤
  </pre>
</div>

<p align="center">
  <strong>The official Go client for Deimos. Simple, efficient, and your reliable link to the foundation of your microservices architecture.</strong>
</p>

---

## Deimos & Phobos: A Symbiotic Relationship

In the cosmos of microservices, **Deimos** and **Phobos** are two celestial bodies orbiting the same planet: your application. They are designed to work in perfect harmony, each fulfilling a critical role.

*   **Deimos (Dread): The Foundation of Knowledge.** Deimos is the distributed, consistent key-value store that acts as the central nervous system for your services. It provides service discovery, configuration management, and distributed coordination.

*   **Phobos (Fear): The Engine of Communication.** Phobos is the RPC framework that governs the interactions between your services, providing speed, resilience, and intelligence for high-performance communication.

**Deimos Client** is the bridge that connects your Go applications to this powerful ecosystem. It provides a simple and idiomatic way to interact with the Deimos cluster, allowing you to leverage its features for service discovery, configuration, and more.

## Features

- **Fluent, Chainable API**: An intuitive and easy-to-use API for all key-value operations.
- **Automatic Cluster Awareness**: Seamlessly handles node discovery and updates, ensuring your client is always connected to a healthy Deimos cluster.
- **Client-Side Load Balancing**: Intelligently distributes requests across all available nodes in the cluster to ensure high availability and performance.
- **Type-Safe by Design**: A fully type-safe interface to minimize runtime errors and improve developer productivity.
- **Highly Extensible**: Easily configured with custom HTTP clients, timeouts, and other options to fit your specific needs.
- **Built for Resilience**: Designed to be fault-tolerant, with built-in mechanisms to handle node failures gracefully.

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
	// Create a new Deimos client, pointing to your cluster endpoints
	client := deimosclient.New(deimosclient.WithEndpoints("http://127.0.0.1:4001", "http://127.0.0.1:4002"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Set a key-value pair
	fmt.Println("Setting value...")
	setResp, err := client.Set(ctx, "/mykey", "hello from the deimos client!")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("-> Set Response: %+v

", setResp)

	// 2. Get the value of the key
	fmt.Println("Getting value...")
	getResp, err := client.Get(ctx, "/mykey")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("-> Get Response: %+v

", getResp)

	// 3. Delete the key
	fmt.Println("Deleting value...")
	delResp, err := client.Delete(ctx, "/mykey")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("-> Delete Response: %+v

", delResp)
}
```

## API Usage Examples

### Creating a Client

```go
// Create a client connected to a single node
client := deimosclient.New(deimosclient.WithEndpoints("http://127.0.0.1:4001"))

// Create a client connected to multiple nodes for high availability
client = deimosclient.New(deimosclient.WithEndpoints("http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"))
```

### Setting Key-Value Pairs

```go
// Set a simple key-value pair
resp, err := client.Set(ctx, "/foo", "bar")

// Set a key-value pair with a TTL (Time-To-Live) of 10 seconds
resp, err = client.Set(ctx, "/foo", "bar", deimosclient.WithTTL(10*time.Second))
```

### Getting Key-Value Pairs

```go
// Get the value of a single key
resp, err := client.Get(ctx, "/foo")

// Recursively get all key-values under a directory
resp, err = client.Get(ctx, "/dir", deimosclient.WithRecursive())
```

### Deleting Key-Value Pairs

```go
// Delete a single key
resp, err := client.Delete(ctx, "/foo")

// Recursively delete a directory and all its contents
resp, err = client.Delete(ctx, "/dir", deimosclient.WithRecursive())
```

### Watching for Changes

The `Watch` feature is a powerful tool for building reactive applications that respond to changes in your Deimos cluster in real-time.

```go
// Watch for changes on a single key
watcher := client.Watch(ctx, "/foo")
for resp := range watcher {
    fmt.Printf("Key '/foo' changed: %+v", resp)
}

// Recursively watch for changes in a directory
watcher = client.Watch(ctx, "/dir", deimosclient.WithRecursive())
for resp := range watcher {
    fmt.Printf("A key in '/dir' changed: %+v", resp)
}
```

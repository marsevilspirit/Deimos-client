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
⠀⠀⠀⠈⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⠀⠉⠉⠉⠀⠀⠀⢸⣿⡟⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⠃⣶⣾⣿⡇⠀⠀⠀⠀
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
- **Distributed Locking**: Robust distributed lock implementation with TTL, auto-renewal, and watch-based coordination.
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

### Distributed Locking

Deimos Client provides a powerful distributed locking mechanism that ensures mutual exclusion across your distributed system. This is essential for coordinating access to shared resources and preventing race conditions.

#### Basic Lock Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    deimos "github.com/marsevilspirit/deimos-client"
)

func main() {
    client := deimos.NewClient([]string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"})
    ctx := context.Background()

    // Create a distributed lock
    lock := client.NewDistributedLock("/locks/my-resource", "node-123")

    // Acquire the lock
    if err := lock.Lock(ctx); err != nil {
        log.Fatalf("Failed to acquire lock: %v", err)
    }
    fmt.Println("✓ Lock acquired successfully")

    // Perform critical section work
    fmt.Println("Performing critical operations...")
    time.Sleep(2 * time.Second)

    // Release the lock
    if err := lock.Unlock(ctx); err != nil {
        log.Fatalf("Failed to release lock: %v", err)
    }
    fmt.Println("✓ Lock released successfully")
}
```

#### Advanced Lock Configuration

```go
// Create a lock with custom TTL and auto-renewal
lock := client.NewDistributedLock("/locks/my-resource", "node-123",
    deimos.WithTTL(30*time.Second),           // Lock expires after 30 seconds
    deimos.WithAutoRenewal(true),             // Enable automatic renewal
    deimos.WithRenewalPeriod(10*time.Second), // Renew every 10 seconds
)

// Acquire lock and start auto-renewal
if err := lock.Lock(ctx); err != nil {
    log.Fatalf("Failed to acquire lock: %v", err)
}

// Start automatic renewal (keeps the lock alive)
lock.StartAutoRenewal(ctx, 10*time.Second)

// Perform long-running work
time.Sleep(60 * time.Second) // Lock will be automatically renewed

// Stop auto-renewal and release lock
lock.StopAutoRenewal()
lock.Unlock(ctx)
```

#### Try Lock (Non-blocking)

```go
// Try to acquire lock without blocking
lock := client.NewDistributedLock("/locks/my-resource", "node-123")

if err := lock.TryLock(ctx); err != nil {
    fmt.Printf("Could not acquire lock immediately: %v\n", err)
    return
}

fmt.Println("Lock acquired immediately!")
defer lock.Unlock(ctx)

// Perform work...
```

#### WithLock Helper Method

```go
// Execute code within a lock using the convenient WithLock method
err := client.WithLock(ctx, "/locks/my-resource", "node-123", func() error {
    fmt.Println("Executing critical section...")
    time.Sleep(2 * time.Second)
    return nil
})

if err != nil {
    log.Fatalf("WithLock failed: %v", err)
}
```

#### Lock Status and Information

```go
lock := client.NewDistributedLock("/locks/my-resource", "node-123")

// Check if lock is currently held
if lock.IsHeld() {
    fmt.Println("Lock is currently held by this client")
}

// Get detailed lock information
info := lock.Info()
fmt.Printf("Lock Info: Key=%s, Held=%v, TTL=%v, LastIndex=%d\n", 
    info.Key, info.Held, info.TTL, info.LastIndex)
```

#### Handling Lock Failures

```go
lock := client.NewDistributedLock("/locks/my-resource", "node-123",
    deimos.WithTTL(10*time.Second))

if err := lock.Lock(ctx); err != nil {
    switch {
    case errors.Is(err, deimos.ErrLockTimeout):
        fmt.Println("Timeout waiting for lock")
    case errors.Is(err, deimos.ErrLockAlreadyHeld):
        fmt.Println("Lock is already held by another client")
    default:
        fmt.Printf("Unexpected error: %v\n", err)
    }
    return
}

defer func() {
    if err := lock.Unlock(ctx); err != nil {
        log.Printf("Failed to release lock: %v", err)
    }
}()

// Critical section...
```

#### Multiple Lock Coordination

```go
// Acquire multiple locks in a specific order to avoid deadlocks
locks := []*deimos.DistributedLock{
    client.NewDistributedLock("/locks/resource-a", "node-123"),
    client.NewDistributedLock("/locks/resource-b", "node-123"),
    client.NewDistributedLock("/locks/resource-c", "node-123"),
}

// Acquire all locks
for i, lock := range locks {
    if err := lock.Lock(ctx); err != nil {
        // Release any previously acquired locks
        for j := i - 1; j >= 0; j-- {
            locks[j].Unlock(ctx)
        }
        log.Fatalf("Failed to acquire lock %d: %v", i, err)
    }
}

// Perform coordinated work with all resources
fmt.Println("All locks acquired, performing coordinated work...")

// Release all locks in reverse order
for i := len(locks) - 1; i >= 0; i-- {
    locks[i].Unlock(ctx)
}
```

## Running Examples

The project includes several examples demonstrating different use cases:

```bash
# Basic key-value operations
go run example/basic/main.go

# Basic distributed lock usage
go run example/lock_basic/main.go

# Complete distributed lock examples
go run example/distributed_lock/main.go

# Advanced lock features
go run example/lock_advanced/main.go

# Multiple client lock watching
go run example/multiple_watch_lock/main.go
```

## Integration Tests

Run the complete integration test suite:

```bash
# Start Deimos cluster and run all tests
./integration-tests.sh
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request or create an Issue.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  <strong>Build better distributed systems, starting with Deimos.</strong>
</p>

## 中文文档

[中文 README](README_zh.md)

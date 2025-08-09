package main

import (
	"context"
	"fmt"
	"time"

	deimos "github.com/marsevilspirit/deimos-client"
)

func main() {
	endpoints := []string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"}
	client := deimos.NewClient(endpoints)

	ctx := context.Background()

	dirKey := "/my-atomic-dir"

	// 1. Create a directory with a 5-second TTL.
	fmt.Printf("1. Creating directory '%s'...\n", dirKey)
	// Make sure the key is clean before we start.
	_, _ = client.Delete(ctx, dirKey, deimos.WithDir(), deimos.WithRecursive())
	createResp, err := client.Set(ctx, dirKey, "", deimos.WithDir())
	if err != nil {
		fmt.Printf("   Failed to create directory: %v\n", err)
		return
	}
	fmt.Printf("   Directory created successfully! Initial index: %d\n", createResp.Node.ModifiedIndex)

	fmt.Println("\n2. Waiting for 3 seconds...")
	time.Sleep(3 * time.Second)

	key := "key"

	fmt.Printf("\n3. Refreshing...\n")
	refreshResp, err := client.Set(ctx, dirKey+"/"+key, "value")
	if err != nil {
		fmt.Printf("   Failed to refresh: %v\n", err)
		return
	}
	fmt.Printf("   refreshed successfully! New index: %d\n", refreshResp.Node.ModifiedIndex)
	fmt.Println("\n--- Demonstration Succeeded ---")
}

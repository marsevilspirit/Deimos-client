package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	deimos "github.com/marsevilspirit/deimos-client"
)

func main() {
	// 确保你的 deimos 服务器正在运行
	endpoints := []string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"}
	client := deimos.NewClient(endpoints)

	fmt.Println("=== Deimos 分布式锁示例 ===\n")

	// 示例1: 基本锁使用
	fmt.Println("1. 基本锁使用示例")
	basicLockExample(client)

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 示例2: 并发锁竞争
	fmt.Println("2. 并发锁竞争示例")
	concurrentLockExample(client)

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 示例3: 锁续约示例
	fmt.Println("3. 锁自动续约示例")
	autoRenewalExample(client)

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 示例4: WithLock 便捷方法
	fmt.Println("4. WithLock 便捷方法示例")
	withLockExample(client)
}

// 基本锁使用示例
func basicLockExample(client *deimos.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	lockKey := "/locks/basic-example"
	nodeID := "node-1"

	// 创建分布式锁
	lock := client.NewDistributedLock(lockKey, nodeID,
		deimos.WithTTL(10*time.Second),
		deimos.WithAutoRenewal(false))

	fmt.Printf("   尝试获取锁: %s\n", lockKey)

	// 获取锁
	if err := lock.Lock(ctx); err != nil {
		log.Printf("   获取锁失败: %v", err)
		return
	}
	fmt.Printf("   ✓ 成功获取锁，持有者: %s\n", nodeID)

	// 检查锁状态
	info := lock.Info()
	fmt.Printf("   锁信息: Key=%s, Held=%t, TTL=%v\n", info.Key, info.Held, info.TTL)

	// 模拟工作
	fmt.Println("   执行受保护的工作...")
	time.Sleep(2 * time.Second)

	// 释放锁
	if err := lock.Unlock(ctx); err != nil {
		log.Printf("   释放锁失败: %v", err)
		return
	}
	fmt.Printf("   ✓ 成功释放锁\n")
}

// 并发锁竞争示例
func concurrentLockExample(client *deimos.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	lockKey := "/locks/concurrent-example"
	numWorkers := 3

	var wg sync.WaitGroup
	results := make(chan string, numWorkers)

	// 启动多个工作者竞争同一个锁
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			nodeID := fmt.Sprintf("worker-%d", workerID)
			lock := client.NewDistributedLock(lockKey, nodeID,
				deimos.WithTTL(5*time.Second),
				deimos.WithAutoRenewal(false))

			fmt.Printf("   Worker %d 尝试获取锁...\n", workerID)

			// 获取锁（内部使用 watchAndLock 机制）
			err := lock.Lock(ctx)
			if err != nil {
				results <- fmt.Sprintf("Worker %d 获取锁失败: %v", workerID, err)
				return
			}

			results <- fmt.Sprintf("✓ Worker %d 成功获取锁", workerID)

			// 模拟工作
			time.Sleep(2 * time.Second)

			// 释放锁
			if err := lock.Unlock(ctx); err != nil {
				results <- fmt.Sprintf("Worker %d 释放锁失败: %v", workerID, err)
			} else {
				results <- fmt.Sprintf("✓ Worker %d 成功释放锁", workerID)
			}
		}(i)
	}

	// 等待所有工作者完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 打印结果
	for result := range results {
		fmt.Printf("   %s\n", result)
	}
}

// 锁自动续约示例
func autoRenewalExample(client *deimos.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	lockKey := "/locks/renewal-example"
	nodeID := "renewal-node"

	// 创建带自动续约的锁
	ttl := 5 * time.Second           // 较短的TTL
	renewalPeriod := 2 * time.Second // 每2秒续约一次
	lock := client.NewDistributedLock(lockKey, nodeID,
		deimos.WithTTL(ttl),
		deimos.WithRenewalPeriod(renewalPeriod),
		deimos.WithAutoRenewal(true))

	fmt.Printf("   获取锁并启动自动续约 (TTL: %v, 续约间隔: %v)\n", ttl, renewalPeriod)

	// 获取锁
	if err := lock.Lock(ctx); err != nil {
		log.Printf("   获取锁失败: %v", err)
		return
	}
	fmt.Printf("   ✓ 成功获取锁\n")

	// 启动自动续约
	lock.StartAutoRenewal(ctx, renewalPeriod)

	// 模拟长时间工作（超过原始TTL）
	fmt.Println("   开始长时间工作（10秒）...")
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		if lock.IsHeld() {
			fmt.Printf("   第 %d 秒: 锁仍然有效\n", i+1)
		} else {
			fmt.Printf("   第 %d 秒: 锁已失效！\n", i+1)
			break
		}
	}

	// 释放锁
	if err := lock.Unlock(ctx); err != nil {
		log.Printf("   释放锁失败: %v", err)
	} else {
		fmt.Printf("   ✓ 成功释放锁\n")
	}
}

// WithLock 便捷方法示例
func withLockExample(client *deimos.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	lockKey := "/locks/withlock-example"
	nodeID := "withlock-node"

	lock := client.NewDistributedLock(lockKey, nodeID,
		deimos.WithTTL(10*time.Second),
		deimos.WithAutoRenewal(false))

	fmt.Printf("   使用 WithLock 方法执行受保护的操作\n")

	// 使用 WithLock 方法，自动处理锁的获取和释放
	err := lock.WithLock(ctx, func() error {
		fmt.Printf("   ✓ 在锁保护下执行操作\n")

		// 模拟一些工作
		for i := 1; i <= 3; i++ {
			fmt.Printf("   执行步骤 %d/3\n", i)
			time.Sleep(1 * time.Second)
		}

		fmt.Printf("   ✓ 操作完成\n")
		return nil
	})

	if err != nil {
		log.Printf("   WithLock 执行失败: %v", err)
	} else {
		fmt.Printf("   ✓ WithLock 执行成功，锁已自动释放\n")
	}
}

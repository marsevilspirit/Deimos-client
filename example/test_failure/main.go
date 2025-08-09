package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	deimos "github.com/marsevilspirit/deimos-client"
)

func main() {
	// 设置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("=== Deimos 失败测试示例 ===")
	fmt.Println("这个测试故意设计为失败，用于验证集成测试脚本的错误处理")
	fmt.Println()

	// 连接到 Deimos 集群
	endpoints := []string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"}
	client := deimos.NewClient(endpoints)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试1: 尝试连接到不存在的端点（这会成功，因为客户端有容错机制）
	fmt.Println("测试1: 基本连接测试...")

	// 测试2: 尝试获取一个锁
	fmt.Println("测试2: 锁操作测试...")
	lockKey := "/locks/test-failure"
	nodeID := "failure-test-node"

	lock := client.NewDistributedLock(lockKey, nodeID,
		deimos.WithTTL(5*time.Second),
		deimos.WithAutoRenewal(false))

	if err := lock.Lock(ctx); err != nil {
		fmt.Printf("✗ 获取锁失败: %v\n", err)
		fmt.Println("测试失败：无法获取锁")
		os.Exit(1)
	}

	fmt.Println("✓ 获取锁成功")

	// 测试3: 故意制造失败 - 尝试操作一个无效的锁
	fmt.Println("测试3: 故意失败测试...")

	// 创建一个会失败的场景：尝试释放一个不存在的锁
	invalidLock := client.NewDistributedLock("/locks/non-existent", "invalid-node")

	// 不先获取锁就尝试释放，这应该会失败
	if err := invalidLock.Unlock(ctx); err != nil {
		fmt.Printf("✗ 预期的失败发生: %v\n", err)
		// 这是我们故意制造的失败
	} else {
		fmt.Println("✗ 意外成功：本应该失败的操作成功了")
	}

	// 释放正常的锁
	if err := lock.Unlock(ctx); err != nil {
		fmt.Printf("✗ 释放锁失败: %v\n", err)
	} else {
		fmt.Println("✓ 锁已释放")
	}

	// 故意让测试失败
	fmt.Println()
	fmt.Println("=== 测试结果 ===")
	fmt.Println("✗ 测试故意失败 - 这是为了验证集成测试脚本的错误处理能力")
	fmt.Println("如果你看到这条消息，说明失败检测机制正在工作")

	// 以非零退出码退出，表示测试失败
	os.Exit(1)
}

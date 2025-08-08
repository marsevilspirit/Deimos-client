package main

import (
	"context"
	"fmt"
	"log"
	"time"

	deimos "github.com/marsevilspirit/deimos-client"
)

func main() {
	// 连接到 Deimos 集群
	endpoints := []string{"http://127.0.0.1:4001"}
	client := deimos.NewClient(endpoints)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建分布式锁
	lockKey := "/locks/my-resource"
	nodeID := "node-123"

	lock := client.NewDistributedLock(lockKey, nodeID) // 使用默认选项

	fmt.Printf("尝试获取分布式锁: %s\n", lockKey)

	// 获取锁
	if err := lock.Lock(ctx); err != nil {
		log.Fatalf("获取锁失败: %v", err)
	}
	fmt.Printf("✓ 成功获取锁，节点ID: %s\n", nodeID)

	// 执行需要互斥访问的操作
	fmt.Println("执行受保护的操作...")
	time.Sleep(5 * time.Second)
	fmt.Println("操作完成")

	// 释放锁
	if err := lock.Unlock(ctx); err != nil {
		log.Fatalf("释放锁失败: %v", err)
	}
	fmt.Println("✓ 成功释放锁")
}

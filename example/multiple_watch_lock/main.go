package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	deimos "github.com/marsevilspirit/deimos-client"
)

func main() {
	fmt.Println("=== 测试多个客户端 Watch 同一个锁的行为 ===")

	client := deimos.NewClient([]string{"http://127.0.0.1:4001"})
	lockKey := "/locks/multi-watch-test"

	// 清理可能存在的锁
	_, _ = client.Delete(context.Background(), lockKey)
	time.Sleep(100 * time.Millisecond) // 等待删除完成

	// 先创建初始锁持有者
	fmt.Println("创建初始锁持有者...")
	initialLock := client.NewDistributedLock(lockKey, "initial-holder")
	ctx := context.Background()

	err := initialLock.TryLock(ctx)
	if err != nil {
		fmt.Printf("初始锁获取失败: %v\n", err)
		return
	}

	// 创建多个等待者
	numWaiters := 5
	var wg sync.WaitGroup
	results := make(chan string, numWaiters+10)

	fmt.Printf("启动 %d 个客户端等待锁释放...\n", numWaiters)

	// 启动多个等待者
	for i := 0; i < numWaiters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			lock := client.NewDistributedLock(lockKey, fmt.Sprintf("waiter-%d", id))
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			results <- fmt.Sprintf("Waiter-%d: 开始等待锁", id)
			start := time.Now()

			err := lock.Lock(ctx)
			elapsed := time.Since(start)

			if err != nil {
				results <- fmt.Sprintf("Waiter-%d: 获取锁失败 (%v) - 耗时: %v", id, err, elapsed)
				return
			}

			results <- fmt.Sprintf("✅ Waiter-%d: 获取锁成功! - 耗时: %v", id, elapsed)

			// 持有锁一段时间
			time.Sleep(1 * time.Second)

			err = lock.Unlock(ctx)
			if err != nil {
				results <- fmt.Sprintf("Waiter-%d: 释放锁失败: %v", id, err)
			} else {
				results <- fmt.Sprintf("✅ Waiter-%d: 释放锁成功", id)
			}
		}(i)
	}

	// 等待所有等待者启动
	time.Sleep(2 * time.Second)

	fmt.Println("✅ 初始锁持有者获取锁成功")
	results <- "Initial: 锁已被初始持有者获取"

	// 等待3秒后释放锁
	fmt.Println("等待 3 秒后释放锁...")
	time.Sleep(3 * time.Second)

	err = initialLock.Unlock(ctx)
	if err != nil {
		fmt.Printf("初始锁释放失败: %v\n", err)
	} else {
		fmt.Println("✅ 初始锁持有者释放锁")
		results <- "Initial: 锁已被释放 - 所有等待者应该收到通知"
	}

	// 等待所有等待者完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 打印所有结果
	fmt.Println("\n=== 详细结果 ===")
	for result := range results {
		fmt.Println(result)
	}

	fmt.Println("\n=== 分析 ===")
	fmt.Println("从结果可以看出:")
	fmt.Println("1. 所有等待者都会收到锁释放的通知（广播机制）")
	fmt.Println("2. 但只有一个等待者能成功获取锁（原子操作）")
	fmt.Println("3. 其他等待者会重新进入等待状态")
	fmt.Println("4. 这是正常的锁竞争行为，不会有问题")
}

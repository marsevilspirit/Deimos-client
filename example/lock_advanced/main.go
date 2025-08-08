package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	deimos "github.com/marsevilspirit/deimos-client"
)

func main() {
	// 设置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 连接到 Deimos 集群
	endpoints := []string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"}
	client := deimos.NewClient(endpoints)

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理优雅关闭
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("=== Deimos 分布式锁高级示例 ===")
	fmt.Println("按 Ctrl+C 优雅退出")
	fmt.Println()

	// 示例1: 锁监控和健康检查
	go lockMonitoringExample(ctx, client)

	// 示例2: 锁超时和重试策略
	go lockTimeoutExample(ctx, client)

	// 示例3: 多个锁的协调
	go multipleLockExample(ctx, client)

	// 等待信号
	<-sigCh
	fmt.Println("\n收到退出信号，正在优雅关闭...")
	cancel()

	// 等待一段时间让所有goroutine清理
	time.Sleep(2 * time.Second)
	fmt.Println("程序已退出")
}

// 锁监控和健康检查示例
func lockMonitoringExample(ctx context.Context, client *deimos.Client) {
	lockKey := "/locks/monitoring-example"
	nodeID := "monitor-node"

	lock := client.NewDistributedLock(lockKey, nodeID,
		deimos.WithTTL(10*time.Second),
		deimos.WithRenewalPeriod(3*time.Second),
		deimos.WithAutoRenewal(true))

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fmt.Printf("[监控] 尝试获取锁: %s\n", lockKey)

		if err := lock.Lock(ctx); err != nil {
			log.Printf("[监控] 获取锁失败: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Printf("[监控] ✓ 获取锁成功\n")

		// 启动自动续约
		lock.StartAutoRenewal(ctx, 3*time.Second)

		// 监控锁状态
		monitorTicker := time.NewTicker(1 * time.Second)
		workDone := make(chan struct{})

		// 模拟工作
		go func() {
			defer close(workDone)
			for i := 0; i < 15; i++ {
				select {
				case <-ctx.Done():
					return
				default:
				}

				fmt.Printf("[监控] 执行工作步骤 %d/15\n", i+1)
				time.Sleep(1 * time.Second)

				// 检查锁状态
				if !lock.IsHeld() {
					fmt.Printf("[监控] ⚠️  锁已丢失！\n")
					return
				}
			}
		}()

		// 监控循环
	monitorLoop:
		for {
			select {
			case <-ctx.Done():
				break monitorLoop
			case <-workDone:
				fmt.Printf("[监控] 工作完成\n")
				break monitorLoop
			case <-monitorTicker.C:
				info := lock.Info()
				if info.Held {
					fmt.Printf("[监控] 锁状态: 正常 (Index: %d)\n", info.LastIndex)
				} else {
					fmt.Printf("[监控] ⚠️  锁状态: 已丢失\n")
					break monitorLoop
				}
			}
		}

		monitorTicker.Stop()

		// 释放锁
		if err := lock.Unlock(ctx); err != nil {
			log.Printf("[监控] 释放锁失败: %v", err)
		} else {
			fmt.Printf("[监控] ✓ 锁已释放\n")
		}

		// 等待一段时间再重新开始
		time.Sleep(5 * time.Second)
	}
}

// 锁超时和重试策略示例
func lockTimeoutExample(ctx context.Context, client *deimos.Client) {
	time.Sleep(2 * time.Second) // 错开启动时间

	lockKey := "/locks/timeout-example"
	nodeID := "timeout-node"

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		lock := client.NewDistributedLock(lockKey, nodeID,
			deimos.WithTTL(5*time.Second),
			deimos.WithAutoRenewal(false))

		fmt.Printf("[超时] 尝试获取锁（使用 watchAndLock 机制）\n")

		// Lock 方法内部已经使用了 watchAndLock，无需手动重试
		err := lock.Lock(ctx)
		if err != nil {
			log.Printf("[超时] 获取锁失败: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf("[超时] ✓ 获取锁成功\n")

		// 模拟工作，但可能超时
		workCtx, workCancel := context.WithTimeout(ctx, 8*time.Second)

		func() {
			defer workCancel()

			for i := 0; i < 10; i++ {
				select {
				case <-workCtx.Done():
					fmt.Printf("[超时] 工作被取消: %v\n", workCtx.Err())
					return
				default:
				}

				fmt.Printf("[超时] 工作进度: %d/10\n", i+1)
				time.Sleep(1 * time.Second)
			}
			fmt.Printf("[超时] 工作完成\n")
		}()

		// 释放锁
		if err := lock.Unlock(ctx); err != nil {
			log.Printf("[超时] 释放锁失败: %v", err)
		} else {
			fmt.Printf("[超时] ✓ 锁已释放\n")
		}

		time.Sleep(3 * time.Second)
	}
}

// 多个锁的协调示例
func multipleLockExample(ctx context.Context, client *deimos.Client) {
	time.Sleep(4 * time.Second) // 错开启动时间

	nodeID := "multi-lock-node"
	lockKeys := []string{
		"/locks/resource-a",
		"/locks/resource-b",
		"/locks/resource-c",
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fmt.Printf("[多锁] 尝试获取多个锁\n")

		// 创建多个锁
		locks := make([]*deimos.DistributedLock, len(lockKeys))
		for i, key := range lockKeys {
			locks[i] = client.NewDistributedLock(key, nodeID,
				deimos.WithTTL(15*time.Second),
				deimos.WithAutoRenewal(false))
		}

		// 尝试按顺序获取所有锁
		acquiredLocks := make([]*deimos.DistributedLock, 0, len(locks))
		success := true

		for i, lock := range locks {
			fmt.Printf("[多锁] 获取锁 %d/%d: %s\n", i+1, len(locks), lockKeys[i])

			if err := lock.TryLock(ctx); err != nil {
				fmt.Printf("[多锁] 获取锁失败: %v\n", err)
				success = false
				break
			}

			acquiredLocks = append(acquiredLocks, lock)
			fmt.Printf("[多锁] ✓ 锁 %d 获取成功\n", i+1)
		}

		if success {
			fmt.Printf("[多锁] ✓ 所有锁获取成功，开始协调工作\n")

			// 执行需要多个资源的工作
			for i := 0; i < 5; i++ {
				fmt.Printf("[多锁] 协调工作进度: %d/5\n", i+1)
				time.Sleep(2 * time.Second)
			}

			fmt.Printf("[多锁] 协调工作完成\n")
		}

		// 释放所有已获取的锁（按相反顺序）
		var wg sync.WaitGroup
		for i := len(acquiredLocks) - 1; i >= 0; i-- {
			wg.Add(1)
			go func(lock *deimos.DistributedLock, index int) {
				defer wg.Done()
				if err := lock.Unlock(ctx); err != nil {
					log.Printf("[多锁] 释放锁 %d 失败: %v", index, err)
				} else {
					fmt.Printf("[多锁] ✓ 锁 %d 已释放\n", index+1)
				}
			}(acquiredLocks[i], i)
		}
		wg.Wait()

		fmt.Printf("[多锁] 所有锁已释放\n")
		time.Sleep(10 * time.Second)
	}
}

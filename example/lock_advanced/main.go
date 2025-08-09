package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	deimos "github.com/marsevilspirit/deimos-client"
)

func main() {
	// 设置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 连接到 Deimos 集群
	endpoints := []string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"}
	client := deimos.NewClient(endpoints)

	// 创建有超时的上下文 - 集成测试模式，60秒后自动退出
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Println("=== Deimos 分布式锁高级集成测试 ===")
	fmt.Println("测试将在60秒后自动结束")
	fmt.Println()

	// 用于收集测试结果
	results := make(chan TestResult, 3)
	var wg sync.WaitGroup

	// 示例1: 锁监控和健康检查
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := lockMonitoringExample(ctx, client)
		results <- result
	}()

	// 示例2: 锁超时和重试策略
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := lockTimeoutExample(ctx, client)
		results <- result
	}()

	// 示例3: 多个锁的协调
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := multipleLockExample(ctx, client)
		results <- result
	}()

	// 等待所有测试完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	testResults := make([]TestResult, 0, 3)
	for result := range results {
		testResults = append(testResults, result)
	}

	// 输出测试总结
	fmt.Println("\n=== 测试结果总结 ===")
	allPassed := true
	for _, result := range testResults {
		status := "✓ PASS"
		if !result.Success {
			status = "✗ FAIL"
			allPassed = false
		}
		fmt.Printf("%s: %s - %s\n", status, result.TestName, result.Message)
	}

	fmt.Printf("\n总体结果: ")
	if allPassed {
		fmt.Println("✓ 所有测试通过")
		os.Exit(0)
	} else {
		fmt.Println("✗ 部分测试失败")
		os.Exit(1)
	}
}

type TestResult struct {
	TestName string
	Success  bool
	Message  string
}

// 锁监控和健康检查示例
func lockMonitoringExample(ctx context.Context, client *deimos.Client) TestResult {
	testName := "锁监控和健康检查"
	lockKey := "/locks/monitoring-example"
	nodeID := "monitor-node"

	fmt.Printf("[监控] 开始测试: %s\n", testName)

	lock := client.NewDistributedLock(lockKey, nodeID,
		deimos.WithTTL(10*time.Second),
		deimos.WithRenewalPeriod(3*time.Second),
		deimos.WithAutoRenewal(true))

	// 尝试获取锁
	fmt.Printf("[监控] 尝试获取锁: %s\n", lockKey)
	if err := lock.Lock(ctx); err != nil {
		return TestResult{
			TestName: testName,
			Success:  false,
			Message:  fmt.Sprintf("获取锁失败: %v", err),
		}
	}

	fmt.Printf("[监控] ✓ 获取锁成功\n")

	// 启动自动续约
	lock.StartAutoRenewal(ctx, 3*time.Second)

	// 监控锁状态并执行工作
	workSteps := 10
	successSteps := 0

	for i := 0; i < workSteps; i++ {
		select {
		case <-ctx.Done():
			goto monitorWorkDone
		default:
		}

		fmt.Printf("[监控] 执行工作步骤 %d/%d\n", i+1, workSteps)
		time.Sleep(1 * time.Second)

		// 检查锁状态
		if !lock.IsHeld() {
			fmt.Printf("[监控] ⚠️  锁已丢失！\n")
			break
		}

		info := lock.Info()
		if info.Held {
			fmt.Printf("[监控] 锁状态: 正常 (Index: %d)\n", info.LastIndex)
			successSteps++
		} else {
			fmt.Printf("[监控] ⚠️  锁状态: 已丢失\n")
			break
		}
	}

monitorWorkDone:
	// 释放锁 - 如果锁已经过期，忽略 "Key not found" 错误
	if err := lock.Unlock(ctx); err != nil {
		if !isKeyNotFoundError(err) {
			return TestResult{
				TestName: testName,
				Success:  false,
				Message:  fmt.Sprintf("释放锁失败: %v", err),
			}
		} else {
			fmt.Printf("[监控] 锁已自动过期，无需手动释放\n")
		}
	}

	fmt.Printf("[监控] ✓ 锁已释放\n")

	success := successSteps >= workSteps/2 // 至少完成一半工作才算成功
	message := fmt.Sprintf("完成 %d/%d 工作步骤", successSteps, workSteps)

	return TestResult{
		TestName: testName,
		Success:  success,
		Message:  message,
	}
}

// 锁超时和重试策略示例
func lockTimeoutExample(ctx context.Context, client *deimos.Client) TestResult {
	testName := "锁超时和重试策略"
	time.Sleep(2 * time.Second) // 错开启动时间

	fmt.Printf("[超时] 开始测试: %s\n", testName)

	lockKey := "/locks/timeout-example"
	nodeID := "timeout-node"

	lock := client.NewDistributedLock(lockKey, nodeID,
		deimos.WithTTL(5*time.Second),
		deimos.WithAutoRenewal(false))

	fmt.Printf("[超时] 尝试获取锁（使用 watchAndLock 机制）\n")

	// Lock 方法内部已经使用了 watchAndLock，无需手动重试
	err := lock.Lock(ctx)
	if err != nil {
		return TestResult{
			TestName: testName,
			Success:  false,
			Message:  fmt.Sprintf("获取锁失败: %v", err),
		}
	}

	fmt.Printf("[超时] ✓ 获取锁成功\n")

	// 模拟工作，测试锁的超时行为
	workSteps := 8
	completedSteps := 0

	for i := 0; i < workSteps; i++ {
		select {
		case <-ctx.Done():
			goto timeoutWorkDone
		default:
		}

		fmt.Printf("[超时] 工作进度: %d/%d\n", i+1, workSteps)
		time.Sleep(1 * time.Second)
		completedSteps++

		// 检查锁是否还有效
		if !lock.IsHeld() {
			fmt.Printf("[超时] 锁已超时失效\n")
			break
		}
	}

timeoutWorkDone:
	// 释放锁 - 如果锁已经过期，忽略 "Key not found" 错误
	unlockErr := lock.Unlock(ctx)
	if unlockErr != nil {
		// 检查是否是锁已过期的错误
		if !isKeyNotFoundError(unlockErr) {
			return TestResult{
				TestName: testName,
				Success:  false,
				Message:  fmt.Sprintf("释放锁失败: %v", unlockErr),
			}
		} else {
			fmt.Printf("[超时] 锁已自动过期，无需手动释放\n")
		}
	}

	fmt.Printf("[超时] ✓ 锁已释放\n")

	// 评估测试结果
	success := completedSteps >= 3 // 至少完成3步工作
	message := fmt.Sprintf("完成 %d/%d 工作步骤", completedSteps, workSteps)

	return TestResult{
		TestName: testName,
		Success:  success,
		Message:  message,
	}
}

// 多个锁的协调示例
func multipleLockExample(ctx context.Context, client *deimos.Client) TestResult {
	testName := "多个锁的协调"
	time.Sleep(4 * time.Second) // 错开启动时间

	fmt.Printf("[多锁] 开始测试: %s\n", testName)

	nodeID := "multi-lock-node"
	lockKeys := []string{
		"/locks/resource-a",
		"/locks/resource-b",
		"/locks/resource-c",
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
	allLocksAcquired := true

	for i, lock := range locks {
		fmt.Printf("[多锁] 获取锁 %d/%d: %s\n", i+1, len(locks), lockKeys[i])

		if err := lock.TryLock(ctx); err != nil {
			fmt.Printf("[多锁] 获取锁失败: %v\n", err)
			allLocksAcquired = false
			break
		}

		acquiredLocks = append(acquiredLocks, lock)
		fmt.Printf("[多锁] ✓ 锁 %d 获取成功\n", i+1)
	}

	var workCompleted bool
	if allLocksAcquired {
		fmt.Printf("[多锁] ✓ 所有锁获取成功，开始协调工作\n")

		// 执行需要多个资源的工作
		workSteps := 5
		for i := 0; i < workSteps; i++ {
			select {
			case <-ctx.Done():
				goto multiLockWorkDone
			default:
			}

			fmt.Printf("[多锁] 协调工作进度: %d/%d\n", i+1, workSteps)
			time.Sleep(1 * time.Second)

			if i == workSteps-1 {
				workCompleted = true
			}
		}

		if workCompleted {
			fmt.Printf("[多锁] 协调工作完成\n")
		}
	}

multiLockWorkDone:
	// 释放所有已获取的锁（按相反顺序）
	var wg sync.WaitGroup
	unlockErrors := 0
	for i := len(acquiredLocks) - 1; i >= 0; i-- {
		wg.Add(1)
		go func(lock *deimos.DistributedLock, index int) {
			defer wg.Done()
			if err := lock.Unlock(ctx); err != nil {
				if !isKeyNotFoundError(err) {
					log.Printf("[多锁] 释放锁 %d 失败: %v", index, err)
					unlockErrors++
				} else {
					fmt.Printf("[多锁] 锁 %d 已自动过期\n", index+1)
				}
			} else {
				fmt.Printf("[多锁] ✓ 锁 %d 已释放\n", index+1)
			}
		}(acquiredLocks[i], i)
	}
	wg.Wait()

	fmt.Printf("[多锁] 所有锁已释放\n")

	// 评估测试结果
	success := allLocksAcquired && workCompleted && unlockErrors == 0
	var message string
	if !allLocksAcquired {
		message = fmt.Sprintf("获取了 %d/%d 个锁", len(acquiredLocks), len(lockKeys))
	} else if !workCompleted {
		message = "获取了所有锁但工作未完成"
	} else if unlockErrors > 0 {
		message = fmt.Sprintf("工作完成但有 %d 个锁释放失败", unlockErrors)
	} else {
		message = "成功获取所有锁、完成工作并释放锁"
	}

	return TestResult{
		TestName: testName,
		Success:  success,
		Message:  message,
	}
}

// 检查错误是否是 "Key not found" 类型的错误
func isKeyNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "Key not found") || strings.Contains(errStr, "HTTP 404")
}

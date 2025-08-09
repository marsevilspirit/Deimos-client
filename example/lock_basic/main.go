package main

import (
	"context"
	"fmt"
	"log"
	"time"

	deimos "github.com/marsevilspirit/deimos-client"
	"github.com/marsevilspirit/deimos-client/example/testutil"
	"github.com/stretchr/testify/assert"
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

	// 创建断言测试器
	t := testutil.NewMockT(false) // false 表示断言失败时不退出程序

	fmt.Println("=== Deimos 分布式锁基础示例 (带断言验证) ===\n")
	fmt.Printf("尝试获取分布式锁: %s\n", lockKey)

	// 断言：初始状态锁应该未被持有
	assert.False(t, lock.IsHeld(), "初始状态锁应该未被持有")

	// 获取锁
	err := lock.Lock(ctx)
	if !assert.NoError(t, err, "获取锁应该成功") {
		log.Fatalf("❌ 获取锁失败: %v", err)
	}
	fmt.Printf("✅ 成功获取锁，节点ID: %s\n", nodeID)

	// 断言：获取锁后应该处于持有状态
	assert.True(t, lock.IsHeld(), "获取锁后应该处于持有状态")

	// 断言：验证锁信息
	info := lock.Info()
	assert.Equal(t, lockKey, info.Key, "锁的键应该与设置的键相同")
	assert.True(t, info.Held, "锁信息应该显示为已持有")
	assert.Greater(t, info.TTL, time.Duration(0), "锁的TTL应该大于0")

	fmt.Printf("📊 锁信息: Key=%s, Held=%t, TTL=%v\n", info.Key, info.Held, info.TTL)

	// 执行需要互斥访问的操作
	fmt.Println("执行受保护的操作...")

	// 模拟工作负载
	startTime := time.Now()
	time.Sleep(2 * time.Second) // 减少等待时间以便测试
	workDuration := time.Since(startTime)

	// 断言：工作时间应该在合理范围内
	assert.GreaterOrEqual(t, workDuration, 2*time.Second, "工作时间应该至少2秒")
	assert.Less(t, workDuration, 5*time.Second, "工作时间应该少于5秒")

	fmt.Printf("✅ 操作完成，耗时: %v\n", workDuration)

	// 断言：在释放前锁仍应该被持有
	assert.True(t, lock.IsHeld(), "释放前锁仍应该被持有")

	// 释放锁
	err = lock.Unlock(ctx)
	if !assert.NoError(t, err, "释放锁应该成功") {
		log.Fatalf("❌ 释放锁失败: %v", err)
	}
	fmt.Println("✅ 成功释放锁")

	// 断言：释放锁后应该不再持有
	assert.False(t, lock.IsHeld(), "释放锁后应该不再持有")

	// 验证锁信息更新
	infoAfterUnlock := lock.Info()
	assert.False(t, infoAfterUnlock.Held, "释放后锁信息应该显示为未持有")

	fmt.Println("\n🎉 分布式锁基础操作测试完成！")
}

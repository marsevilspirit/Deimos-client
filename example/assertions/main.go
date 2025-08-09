package main

import (
	"context"
	"fmt"
	"time"

	deimos "github.com/marsevilspirit/deimos-client"
	"github.com/marsevilspirit/deimos-client/example/testutil"
	"github.com/stretchr/testify/assert"
)

func main() {
	fmt.Println("=== Deimos 断言验证示例 ===\n")

	endpoints := []string{"http://127.0.0.1:4001"}
	client := deimos.NewClient(endpoints)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 测试计数器
	testCount := 0
	passedCount := 0

	// 创建断言测试器
	t := testutil.NewMockT(false) // false 表示断言失败时不退出程序

	// 辅助函数：运行测试用例
	runTest := func(name string, testFunc func() bool) {
		testCount++
		fmt.Printf("🧪 测试 %d: %s\n", testCount, name)

		if testFunc() {
			passedCount++
			fmt.Println("   ✅ 通过\n")
		} else {
			fmt.Println("   ❌ 失败\n")
		}
	}

	// 测试1: 基本键值操作断言
	runTest("基本键值操作断言", func() bool {
		key := "/test/basic"
		value := "test-value"

		// 设置键值
		setResp, err := client.Set(ctx, key, value)
		if !assert.NoError(t, err, "设置键值应该成功") {
			return false
		}

		// 断言设置响应
		success := assert.Equal(t, "set", setResp.Action, "操作类型应该是 'set'") &&
			assert.Equal(t, value, setResp.Node.Value, "节点值应该与设置值相同") &&
			assert.Greater(t, setResp.Node.ModifiedIndex, uint64(0), "修改索引应该大于0")

		// 获取键值
		getResp, err := client.Get(ctx, key)
		if !assert.NoError(t, err, "获取键值应该成功") {
			return false
		}

		// 断言获取响应
		success = success &&
			assert.Equal(t, value, getResp.Node.Value, "获取的值应该与设置的值相同") &&
			assert.Equal(t, key, getResp.Node.Key, "获取的键应该与请求的键相同")

		// 清理
		_, _ = client.Delete(ctx, key)

		return success
	})

	// 测试2: 分布式锁状态断言
	runTest("分布式锁状态断言", func() bool {
		lockKey := "/test/lock"
		nodeID := "test-node"

		lock := client.NewDistributedLock(lockKey, nodeID)

		// 断言初始状态
		if !assert.False(t, lock.IsHeld(), "初始状态锁应该未被持有") {
			return false
		}

		// 获取锁
		err := lock.Lock(ctx)
		if !assert.NoError(t, err, "获取锁应该成功") {
			return false
		}

		// 断言锁定状态
		success := assert.True(t, lock.IsHeld(), "获取锁后应该处于持有状态")

		// 验证锁信息
		info := lock.Info()
		success = success &&
			assert.Equal(t, lockKey, info.Key, "锁的键应该与设置的键相同") &&
			assert.True(t, info.Held, "锁信息应该显示为已持有") &&
			assert.Greater(t, info.TTL, time.Duration(0), "锁的TTL应该大于0")

		// 释放锁
		err = lock.Unlock(ctx)
		if !assert.NoError(t, err, "释放锁应该成功") {
			return false
		}

		// 断言释放状态
		success = success && assert.False(t, lock.IsHeld(), "释放锁后应该不再持有")

		return success
	})

	// 测试3: 错误处理断言
	runTest("错误处理断言", func() bool {
		nonExistentKey := "/test/non-existent"

		// 尝试获取不存在的键
		_, err := client.Get(ctx, nonExistentKey)

		// 断言应该返回错误
		return assert.Error(t, err, "获取不存在的键应该返回错误")
	})

	// 测试4: 时间相关断言
	runTest("时间相关断言", func() bool {
		key := "/test/timing"
		value := "timing-test"

		// 测量设置操作时间
		startTime := time.Now()
		_, err := client.Set(ctx, key, value)
		duration := time.Since(startTime)

		if !assert.NoError(t, err, "设置操作应该成功") {
			return false
		}

		// 断言操作时间应该在合理范围内
		success := assert.Less(t, duration, 5*time.Second, "设置操作应该在5秒内完成") &&
			assert.Greater(t, duration, time.Duration(0), "操作时间应该大于0")

		// 清理
		_, _ = client.Delete(ctx, key)

		return success
	})

	// 测试5: 数据类型断言
	runTest("数据类型断言", func() bool {
		key := "/test/types"

		// 测试不同类型的值
		testValues := []string{
			"string-value",
			"123",
			"true",
			"",
			"special-chars-!@#$%^&*()",
		}

		success := true
		for i, value := range testValues {
			testKey := fmt.Sprintf("%s/%d", key, i)

			// 设置值
			setResp, err := client.Set(ctx, testKey, value)
			if !assert.NoError(t, err, fmt.Sprintf("设置值 '%s' 应该成功", value)) {
				success = false
				continue
			}

			// 验证值类型
			success = success &&
				assert.IsType(t, "", setResp.Node.Value, "节点值应该是字符串类型") &&
				assert.Equal(t, value, setResp.Node.Value, fmt.Sprintf("值应该与设置的 '%s' 相同", value))

			// 清理
			_, _ = client.Delete(ctx, testKey)
		}

		return success
	})

	// 输出测试结果统计
	fmt.Printf("📊 测试结果统计:\n")
	fmt.Printf("   总测试数: %d\n", testCount)
	fmt.Printf("   通过数: %d\n", passedCount)
	fmt.Printf("   失败数: %d\n", testCount-passedCount)
	fmt.Printf("   通过率: %.1f%%\n", float64(passedCount)/float64(testCount)*100)

	if passedCount == testCount {
		fmt.Println("\n🎉 所有断言测试通过！")
	} else {
		fmt.Println("\n⚠️  部分测试失败，请检查 Deimos 服务器状态")
	}
}

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
	// 确保你的 deimos 服务器正在运行
	endpoints := []string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"}
	client := deimos.NewClient(endpoints)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	key := "/my-app/session/user123"
	value := "some-session-data"
	ttlDuration := 5 * time.Second

	fmt.Println("=== Deimos TTL (生存时间) 示例 (带断言验证) ===\n")

	// 1. 设置一个带 TTL 的键
	fmt.Printf("1. 设置键: '%s' -> 值: '%s' (TTL: %s)\n", key, value, ttlDuration)

	setStartTime := time.Now()
	setResp, err := client.Set(ctx, key, value, deimos.WithTTL(ttlDuration))
	setDuration := time.Since(setStartTime)

	// 创建断言测试器
	t := testutil.NewMockT(false) // false 表示断言失败时不退出程序

	if !assert.NoError(t, err, "设置带TTL的键应该成功") {
		fmt.Printf("   ❌ 设置失败: %v\n", err)
		return
	}

	// 断言：验证设置操作的响应
	assert.Equal(t, "set", setResp.Action, "操作类型应该是 'set'")
	assert.NotEmpty(t, setResp.Node.Value, "节点值不应该为空")
	assert.Greater(t, setResp.Node.ModifiedIndex, uint64(0), "修改索引应该大于0")
	assert.Less(t, setDuration, 1*time.Second, "设置操作应该在1秒内完成")

	fmt.Printf("   ✅ 设置成功! Action: %s, ModifiedIndex: %d, 耗时: %v\n\n",
		setResp.Action, setResp.Node.ModifiedIndex, setDuration)

	// 2. 立即获取这个键，应该是成功的
	fmt.Printf("2. 立即获取键: '%s'\n", key)

	getStartTime := time.Now()
	getResp, err := client.Get(ctx, key)
	getDuration := time.Since(getStartTime)

	if !assert.NoError(t, err, "在TTL过期前获取键应该成功") {
		fmt.Printf("   ❌ 获取失败: %v\n", err)
		return
	}

	// 断言：验证获取到的值
	assert.Equal(t, value, getResp.Node.Value, "获取的值应该与设置的值相同")
	assert.Equal(t, key, getResp.Node.Key, "获取的键应该与请求的键相同")
	assert.Less(t, getDuration, 1*time.Second, "获取操作应该在1秒内完成")

	fmt.Printf("   ✅ 获取成功! Value: '%s', 耗时: %v\n\n", getResp.Node.Value, getDuration)

	// 3. 等待超过 TTL 的时间
	waitDuration := 6 * time.Second
	fmt.Printf("3. 等待 %s (超过 TTL 时间)...\n", waitDuration)

	waitStartTime := time.Now()
	time.Sleep(waitDuration)
	actualWaitDuration := time.Since(waitStartTime)

	// 断言：验证等待时间
	assert.GreaterOrEqual(t, actualWaitDuration, waitDuration, "实际等待时间应该不少于预期时间")
	assert.Less(t, actualWaitDuration, waitDuration+1*time.Second, "实际等待时间不应该超出太多")

	fmt.Printf("   ⏰ 等待完成，实际耗时: %v\n\n", actualWaitDuration)

	// 4. 再次获取这个键，预期会失败
	fmt.Printf("4. 在 TTL 过期后再次获取键: '%s'\n", key)

	expiredGetStartTime := time.Now()
	_, err = client.Get(ctx, key)
	expiredGetDuration := time.Since(expiredGetStartTime)

	// 断言：验证获取过期键应该失败
	assert.Error(t, err, "获取过期的键应该返回错误")
	assert.Less(t, expiredGetDuration, 2*time.Second, "获取过期键的操作应该快速失败")

	if err != nil {
		fmt.Printf("   ✅ 获取失败 (符合预期): %v, 耗时: %v\n", err, expiredGetDuration)
	} else {
		fmt.Println("   ❌ 错误：键没有按预期过期！")
	}

	// 5. 验证总体时间逻辑
	totalElapsed := time.Since(setStartTime)
	expectedMinTime := ttlDuration + waitDuration

	assert.GreaterOrEqual(t, totalElapsed, expectedMinTime,
		"总耗时应该至少等于TTL时间加等待时间")

	fmt.Printf("\n📊 时间统计:\n")
	fmt.Printf("   - TTL 设置: %v\n", ttlDuration)
	fmt.Printf("   - 等待时间: %v\n", waitDuration)
	fmt.Printf("   - 总耗时: %v\n", totalElapsed)
	fmt.Printf("   - 预期最小时间: %v\n", expectedMinTime)

	fmt.Println("\n🎉 TTL 功能测试完成！")
}

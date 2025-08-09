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
	endpoints := []string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"}
	client := deimos.NewClient(endpoints)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := "/key"
	value := "localhost:8080"

	// 创建断言测试器
	t := testutil.NewMockT(false) // false 表示断言失败时不退出程序

	fmt.Println("=== Deimos 基本操作示例 (带断言验证) ===\n")

	// 1. 设置一个值
	fmt.Printf("1. 设置键: '%s' -> 值: '%s'\n", key, value)
	setResp, err := client.Set(ctx, key, value)
	if !assert.NoError(t, err, "设置键值应该成功") {
		fmt.Printf("   ❌ 设置失败: %+v\n", err)
		return
	}

	// 断言：验证设置操作的响应
	assert.Equal(t, "set", setResp.Action, "操作类型应该是 'set'")
	assert.NotEmpty(t, setResp.Node.Value, "节点值不应该为空")
	assert.Greater(t, setResp.Node.ModifiedIndex, uint64(0), "修改索引应该大于0")

	fmt.Printf("   ✅ 设置成功! Action: %s, ModifiedIndex: %d\n\n", setResp.Action, setResp.Node.ModifiedIndex)

	// 2. 获取这个值
	fmt.Printf("2. 获取键: '%s'\n", key)
	getResp, err := client.Get(ctx, key)
	if !assert.NoError(t, err, "获取键值应该成功") {
		fmt.Printf("   ❌ 获取失败: %v\n", err)
		return
	}

	// 断言：验证获取到的值
	assert.Equal(t, value, getResp.Node.Value, "获取的值应该与设置的值相同")
	assert.Equal(t, key, getResp.Node.Key, "获取的键应该与请求的键相同")
	assert.Greater(t, getResp.Node.CreatedIndex, uint64(0), "创建索引应该大于0")

	fmt.Printf("   ✅ 获取成功! Value: '%s', CreatedIndex: %d\n\n", getResp.Node.Value, getResp.Node.CreatedIndex)

	// 3. 删除这个值
	fmt.Printf("3. 删除键: '%s'\n", key)
	delResp, err := client.Delete(ctx, key)
	if !assert.NoError(t, err, "删除键值应该成功") {
		fmt.Printf("   ❌ 删除失败: %v\n", err)
		return
	}

	// 断言：验证删除操作的响应
	assert.Equal(t, "delete", delResp.Action, "操作类型应该是 'delete'")
	assert.Equal(t, value, delResp.PrevNode.Value, "删除前的值应该与原值相同")
	assert.Equal(t, key, delResp.PrevNode.Key, "删除的键应该与请求的键相同")

	fmt.Printf("   ✅ 删除成功! Action: %s, PrevValue: '%s'\n\n", delResp.Action, delResp.PrevNode.Value)

	// 4. 再次尝试获取，预期会失败
	fmt.Printf("4. 再次获取已删除的键: '%s'\n", key)
	_, err = client.Get(ctx, key)

	// 断言：验证获取已删除的键应该失败
	assert.Error(t, err, "获取已删除的键应该返回错误")

	if err != nil {
		fmt.Printf("   ✅ 获取失败 (符合预期): %v\n", err)
	} else {
		fmt.Println("   ❌ 错误：竟然获取到了已删除的键！")
	}

	fmt.Println("\n🎉 所有基本操作测试完成！")
}

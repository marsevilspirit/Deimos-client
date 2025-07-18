package main

import (
	"context"
	"fmt"
	"time"

	deimos "github.com/marsevilspirit/deimos-client"
)

func main() {
	// 确保你的 etcd 服务器正在运行
	endpoints := []string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"}
	client := deimos.NewClient(endpoints)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	key := "/my-app/session/user123"
	value := "some-session-data"
	ttlDuration := 5 * time.Second

	// 1. 设置一个带 TTL 的键
	fmt.Printf("1. 设置键: '%s' -> 值: '%s' (TTL: %s)\n", key, value, ttlDuration)
	// 使用新的 WithTTL 选项
	setResp, err := client.Set(ctx, key, value, deimos.WithTTL(ttlDuration))
	if err != nil {
		fmt.Printf("   设置失败: %v\n", err)
		return
	}
	fmt.Printf("   设置成功! Action: %s, ModifiedIndex: %d\n\n", setResp.Action, setResp.Node.ModifiedIndex)

	// 2. 立即获取这个键，应该是成功的
	fmt.Printf("2. 立即获取键: '%s'\n", key)
	getResp, err := client.Get(ctx, key)
	if err != nil {
		fmt.Printf("   获取失败: %v\n", err)
		return
	}
	fmt.Printf("   获取成功! Value: '%s'\n\n", getResp.Node.Value)

	// 3. 等待超过 TTL 的时间
	waitDuration := 6 * time.Second
	fmt.Printf("3. 等待 %s (超过 TTL 时间)...\n\n", waitDuration)
	time.Sleep(waitDuration)

	// 4. 再次获取这个键，预期会失败
	fmt.Printf("4. 在 TTL 过期后再次获取键: '%s'\n", key)
	_, err = client.Get(ctx, key)
	if err != nil {
		fmt.Printf("   获取失败 (预期之中): %v\n", err)
	} else {
		fmt.Println("   错误：键没有按预期过期！")
	}
}

package main

import (
	"context"
	"fmt"
	"time"

	deimos "github.com/marsevilspirit/deimos-client"
)

func main() {
	endpoints := []string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"}
	client := deimos.NewClient(endpoints)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := "/key"
	value := "localhost:8080"

	// 1. 设置一个值
	fmt.Printf("1. 设置键: '%s' -> 值: '%s'\n", key, value)
	setResp, err := client.Set(ctx, key, value)
	if err != nil {
		fmt.Printf("   设置失败: %+v\n", err)
		return
	}
	fmt.Printf("   设置成功! Action: %s, ModifiedIndex: %d\n\n", setResp.Action, setResp.Node.ModifiedIndex)

	// 2. 获取这个值
	fmt.Printf("2. 获取键: '%s'\n", key)
	getResp, err := client.Get(ctx, key)
	if err != nil {
		fmt.Printf("   获取失败: %v\n", err)
		return
	}
	fmt.Printf("   获取成功! Value: '%s', CreatedIndex: %d\n\n", getResp.Node.Value, getResp.Node.CreatedIndex)

	// 3. 删除这个值
	fmt.Printf("3. 删除键: '%s'\n", key)
	delResp, err := client.Delete(ctx, key)
	if err != nil {
		fmt.Printf("   删除失败: %v\n", err)
		return
	}
	fmt.Printf("   删除成功! Action: %s, PrevValue: '%s'\n\n", delResp.Action, delResp.PrevNode.Value)

	// 4. 再次尝试获取，预期会失败
	fmt.Printf("4. 再次获取已删除的键: '%s'\n", key)
	_, err = client.Get(ctx, key)
	if err != nil {
		fmt.Printf("   获取失败 (预期之中): %v\n", err)
	} else {
		fmt.Println("   错误：竟然获取到了已删除的键！")
	}
}

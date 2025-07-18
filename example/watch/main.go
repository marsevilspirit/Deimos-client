package main

import (
	"context"
	"fmt"
	"time"

	deimos "github.com/marsevilspirit/deimos-client"
)

// main 函数已更新，用于演示 Watch 功能
func main() {
	endpoints := []string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"}
	client := deimos.NewClient(endpoints)
	// 创建一个可以被取消的 context，用于控制 watcher 的生命周期
	ctx, cancel := context.WithCancel(context.Background())
	// 在 main 函数结束时调用 cancel()，确保所有 goroutine 都能退出。
	defer cancel()

	key := "/foo"

	// 启动一个 goroutine 来监视键的变化
	fmt.Printf(">>> 启动对 '%s' 的监视...\n", key)
	watchChan := client.Watch(ctx, key)

	// 启动另一个 goroutine 来消费 watchChan 的数据
	go func() {
		// for range 会一直阻塞，直到 channel 中有数据或 channel 被关闭
		for resp := range watchChan {
			fmt.Printf("<<< 监听到变化! Action: '%s', Key: '%s', Value: '%s'\n", resp.Action, resp.Node.Key, resp.Node.Value)
		}
		fmt.Println("!!! Watcher channel 已关闭，监听结束。")
	}()

	// 在主 goroutine 中模拟操作
	fmt.Println("... 等待 2 秒后设置一个值 ...")
	time.Sleep(2 * time.Second)
	_, err := client.Set(context.Background(), key, "bar")
	if err != nil {
		panic(err)
	}

	fmt.Println("... 等待 2 秒后更新这个值 ...")
	time.Sleep(2 * time.Second)
	_, err = client.Set(context.Background(), key, "bar2")
	if err != nil {
		panic(err)
	}

	fmt.Println("... 等待 2 秒后删除这个值 ...")
	time.Sleep(2 * time.Second)
	_, err = client.Delete(context.Background(), key)
	if err != nil {
		panic(err)
	}

	// 再等待一会，然后取消 context 来停止 watcher
	fmt.Println("... 等待 2 秒后取消监视 ...")
	time.Sleep(2 * time.Second)
	cancel()

	// 等待一小会，让 "监听结束" 的消息打印出来
	time.Sleep(1 * time.Second)
	fmt.Println(">>> 程序退出。")
}

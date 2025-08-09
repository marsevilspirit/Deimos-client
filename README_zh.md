# Deimos Client: 您的 Go 应用通往 Deimos 生态系统的网关

<div align="center">
  <pre>
⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⠇⠀⠘⢿⣿⣿⣧⣶⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣤⣤⣾⣿⣿⣿⣿⣿⣿⣿⣿⡿⣿⣿⣿⣿⣿⣿⣿⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠉⢿⣿⣿⣿⠿⠟⠋⠉⠀⠀⢸⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠸⣿⣿⣷⣶⣤⣤⣀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣧⣤⣤⣤⡀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠈⠙⠛⠛⠛⢿⣿⣧⠀⢸⣿⣿⣿⣿⣿⣿⣿⠟⠛⠻⣿⣦⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⣴⣶⣶⡇⠀⠀⠀⣿⣿⠀⢼⣿⣿⣿⣿⣿⣿⣿⣤⡀⠀⢻⣿⡇⠀⠀⠀⠀⠀
⠀⠀⠀⠈⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⠀⠉⠉⠉⠀⠀⠀⢸⣿⡟⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⠃⣶⣾⣿⡇⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠈⠉⠙⢻⣿⣿⣿⣿⣿⣿⣿⠀⠀⢸⣿⣿⠇⠀⢸⣿⡇⠀⢸⣿⣿⣿⣿⣿⣿⣿⡏⠀⠛⣿⣿⠇⠀⠀⠀⠀⠀
⣤⣤⣤⣤⣤⣤⣤⣤⣼⣿⣿⣿⣿⣿⣿⣿⣷⣤⣤⣤⣤⣤⣤⣼⣿⣧⣤⣤⣿⣿⣿⣿⣿⣿⣿⣧⣤⣴⣿⣿⣤⣤⣤⣤⣤⣤
  </pre>
</div>

<p align="center">
  <strong>Deimos 官方 Go 客户端。简单、高效，是您连接微服务架构基础的可靠纽带。</strong>
</p>

---

## Deimos 与 Phobos：共生关系

在微服务的宇宙中，**Deimos** 和 **Phobos** 是围绕同一颗行星（您的应用程序）运行的两个天体。它们被设计为完美协调工作，各自履行关键职责。

*   **Deimos（恐惧）：知识的基础。** Deimos 是分布式、一致性键值存储，充当您服务的中央神经系统。它提供服务发现、配置管理和分布式协调。

*   **Phobos（恐慌）：通信的引擎。** Phobos 是管理服务间交互的 RPC 框架，为高性能通信提供速度、弹性和智能。

**Deimos Client** 是连接您的 Go 应用程序到这个强大生态系统的桥梁。它提供了一种简单而惯用的方式来与 Deimos 集群交互，让您能够利用其服务发现、配置管理等功能。

## 特性

- **流畅的链式 API**：直观易用的 API，支持所有键值操作。
- **分布式锁**：强大的分布式锁实现，支持 TTL、自动续约和基于监听的协调。
- **自动集群感知**：无缝处理节点发现和更新，确保您的客户端始终连接到健康的 Deimos 集群。
- **客户端负载均衡**：智能地在集群中的所有可用节点间分发请求，确保高可用性和性能。
- **类型安全设计**：完全类型安全的接口，最小化运行时错误，提高开发者生产力。
- **高度可扩展**：易于配置自定义 HTTP 客户端、超时和其他选项，以满足您的特定需求。
- **为弹性而构建**：设计为容错，内置机制优雅处理节点故障。

## 安装

```bash
go get github.com/marsevilspirit/deimos-client
```

## 快速开始

以下是如何使用 `Deimos-client` 连接到 Deimos 集群并执行基本操作的简单示例。

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	deimosclient "github.com/marsevilspirit/deimos-client"
)

func main() {
	// 创建新的 Deimos 客户端，指向您的集群端点
	client := deimosclient.New(deimosclient.WithEndpoints("http://127.0.0.1:4001", "http://127.0.0.1:4002"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. 设置键值对
	fmt.Println("设置值...")
	setResp, err := client.Set(ctx, "/mykey", "来自 deimos 客户端的问候!")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("-> 设置响应: %+v\n\n", setResp)

	// 2. 获取键的值
	fmt.Println("获取值...")
	getResp, err := client.Get(ctx, "/mykey")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("-> 获取响应: %+v\n\n", getResp)

	// 3. 删除键
	fmt.Println("删除值...")
	delResp, err := client.Delete(ctx, "/mykey")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("-> 删除响应: %+v\n\n", delResp)
}
```

## API 使用示例

### 创建客户端

```go
// 创建连接到单个节点的客户端
client := deimosclient.New(deimosclient.WithEndpoints("http://127.0.0.1:4001"))

// 创建连接到多个节点的客户端以实现高可用性
client = deimosclient.New(deimosclient.WithEndpoints("http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"))
```

### 设置键值对

```go
// 设置简单的键值对
resp, err := client.Set(ctx, "/foo", "bar")

// 设置带有 TTL（生存时间）10 秒的键值对
resp, err = client.Set(ctx, "/foo", "bar", deimosclient.WithTTL(10*time.Second))
```

### 获取键值对

```go
// 获取单个键的值
resp, err := client.Get(ctx, "/foo")

// 递归获取目录下的所有键值
resp, err = client.Get(ctx, "/dir", deimosclient.WithRecursive())
```

### 删除键值对

```go
// 删除单个键
resp, err := client.Delete(ctx, "/foo")

// 递归删除目录及其所有内容
resp, err = client.Delete(ctx, "/dir", deimosclient.WithRecursive())
```

### 监听变化

`Watch` 功能是构建响应式应用程序的强大工具，可以实时响应 Deimos 集群中的变化。

```go
// 监听单个键的变化
watcher := client.Watch(ctx, "/foo")
for resp := range watcher {
    fmt.Printf("键 '/foo' 发生变化: %+v", resp)
}

// 递归监听目录中的变化
watcher = client.Watch(ctx, "/dir", deimosclient.WithRecursive())
for resp := range watcher {
    fmt.Printf("目录 '/dir' 中的键发生变化: %+v", resp)
}
```

### 分布式锁

Deimos Client 提供了强大的分布式锁机制，确保分布式系统中的互斥访问。这对于协调共享资源访问和防止竞态条件至关重要。

#### 基本锁使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    deimos "github.com/marsevilspirit/deimos-client"
)

func main() {
    client := deimos.NewClient([]string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"})
    ctx := context.Background()

    // 创建分布式锁
    lock := client.NewDistributedLock("/locks/my-resource", "node-123")

    // 获取锁
    if err := lock.Lock(ctx); err != nil {
        log.Fatalf("获取锁失败: %v", err)
    }
    fmt.Println("✓ 锁获取成功")

    // 执行临界区工作
    fmt.Println("执行关键操作...")
    time.Sleep(2 * time.Second)

    // 释放锁
    if err := lock.Unlock(ctx); err != nil {
        log.Fatalf("释放锁失败: %v", err)
    }
    fmt.Println("✓ 锁释放成功")
}
```

#### 高级锁配置

```go
// 创建带有自定义 TTL 和自动续约的锁
lock := client.NewDistributedLock("/locks/my-resource", "node-123",
    deimos.WithTTL(30*time.Second),           // 锁在 30 秒后过期
    deimos.WithAutoRenewal(true),             // 启用自动续约
    deimos.WithRenewalPeriod(10*time.Second), // 每 10 秒续约一次
)

// 获取锁并启动自动续约
if err := lock.Lock(ctx); err != nil {
    log.Fatalf("获取锁失败: %v", err)
}

// 启动自动续约（保持锁活跃）
lock.StartAutoRenewal(ctx, 10*time.Second)

// 执行长时间运行的工作
time.Sleep(60 * time.Second) // 锁将自动续约

// 停止自动续约并释放锁
lock.StopAutoRenewal()
lock.Unlock(ctx)
```

#### 尝试锁（非阻塞）

```go
// 尝试获取锁而不阻塞
lock := client.NewDistributedLock("/locks/my-resource", "node-123")

if err := lock.TryLock(ctx); err != nil {
    fmt.Printf("无法立即获取锁: %v\n", err)
    return
}

fmt.Println("锁立即获取成功!")
defer lock.Unlock(ctx)

// 执行工作...
```

#### WithLock 辅助方法

```go
// 使用便捷的 WithLock 方法在锁内执行代码
err := client.WithLock(ctx, "/locks/my-resource", "node-123", func() error {
    fmt.Println("执行临界区...")
    time.Sleep(2 * time.Second)
    return nil
})

if err != nil {
    log.Fatalf("WithLock 失败: %v", err)
}
```

#### 锁状态和信息

```go
lock := client.NewDistributedLock("/locks/my-resource", "node-123")

// 检查锁是否当前被持有
if lock.IsHeld() {
    fmt.Println("锁当前被此客户端持有")
}

// 获取详细的锁信息
info := lock.Info()
fmt.Printf("锁信息: Key=%s, Held=%v, TTL=%v, LastIndex=%d\n", 
    info.Key, info.Held, info.TTL, info.LastIndex)
```

#### 处理锁失败

```go
lock := client.NewDistributedLock("/locks/my-resource", "node-123",
    deimos.WithTTL(10*time.Second))

if err := lock.Lock(ctx); err != nil {
    switch {
    case errors.Is(err, deimos.ErrLockTimeout):
        fmt.Println("等待锁超时")
    case errors.Is(err, deimos.ErrLockAlreadyHeld):
        fmt.Println("锁已被其他客户端持有")
    default:
        fmt.Printf("意外错误: %v\n", err)
    }
    return
}

defer func() {
    if err := lock.Unlock(ctx); err != nil {
        log.Printf("释放锁失败: %v", err)
    }
}()

// 临界区...
```

#### 多锁协调

```go
// 按特定顺序获取多个锁以避免死锁
locks := []*deimos.DistributedLock{
    client.NewDistributedLock("/locks/resource-a", "node-123"),
    client.NewDistributedLock("/locks/resource-b", "node-123"),
    client.NewDistributedLock("/locks/resource-c", "node-123"),
}

// 获取所有锁
for i, lock := range locks {
    if err := lock.Lock(ctx); err != nil {
        // 释放任何之前获取的锁
        for j := i - 1; j >= 0; j-- {
            locks[j].Unlock(ctx)
        }
        log.Fatalf("获取锁 %d 失败: %v", i, err)
    }
}

// 使用所有资源执行协调工作
fmt.Println("所有锁已获取，执行协调工作...")

// 按相反顺序释放所有锁
for i := len(locks) - 1; i >= 0; i-- {
    locks[i].Unlock(ctx)
}
```

## 运行示例

项目包含多个示例，展示不同的使用场景：

```bash
# 基本键值操作
go run example/basic/main.go

# 分布式锁基础用法
go run example/lock_basic/main.go

# 分布式锁完整示例
go run example/distributed_lock/main.go

# 分布式锁高级功能
go run example/lock_advanced/main.go

# 多客户端锁监听
go run example/multiple_watch_lock/main.go
```

## 集成测试

运行完整的集成测试套件：

```bash
# 启动 Deimos 集群并运行所有测试
./integration-tests.sh
```

## 贡献

欢迎贡献！请随时提交 Pull Request 或创建 Issue。

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

---

<p align="center">
  <strong>构建更好的分布式系统，从 Deimos 开始。</strong>
</p>

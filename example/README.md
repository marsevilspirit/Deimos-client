# Deimos Client 示例

本目录包含了使用 Deimos Client 的各种示例，包括基本操作、分布式锁、监听等功能。

## 模块结构

本示例目录使用独立的 Go 模块管理，与主项目分离：

```
example/
├── go.mod              # 独立的模块文件
├── go.sum              # 依赖锁定文件
├── README.md           # 本文件
├── basic/              # 基本操作示例
├── lock_basic/         # 基础锁示例
├── distributed_lock/   # 完整分布式锁示例
├── lock_advanced/      # 高级锁用法
├── watch/              # 监听示例
├── ttl/                # TTL 示例
├── dir/                # 目录操作示例
├── test/               # 测试示例
└── multiple_watch_lock/ # 多重监听锁示例
```

## 运行示例

### 方式一：使用 Makefile（推荐）

```bash
# 从项目根目录运行
make examples
```

### 方式二：直接运行

```bash
# 进入 example 目录
cd example

# 运行特定示例
go run basic/main.go
go run lock_basic/main.go
go run distributed_lock/main.go
```

## 前置条件

确保 Deimos 服务器正在运行：

```bash
# 启动单节点 Deimos 服务器
./bin/deimos --name node1 \
  --listen-client-urls http://127.0.0.1:4001 \
  --advertise-client-urls http://127.0.0.1:4001 \
  --listen-peer-urls http://127.0.0.1:7001 \
  --advertise-peer-urls http://127.0.0.1:7001 \
  --bootstrap-config "node1=http://localhost:7001"
```

## 示例说明

### 1. 基础锁示例 (`lock_basic/`)

最简单的分布式锁使用示例，展示：
- 创建分布式锁
- 获取锁
- 执行受保护的操作
- 释放锁

```bash
cd example/lock_basic
go run main.go
```

### 2. 完整功能示例 (`distributed_lock/`)

展示分布式锁的完整功能，包括：
- 基本锁操作
- 并发锁竞争
- 自动续约机制
- WithLock 便捷方法

```bash
cd example/distributed_lock
go run main.go
```

### 3. 高级用法示例 (`lock_advanced/`)

展示高级用法和最佳实践：
- 锁监控和健康检查
- 超时和重试策略
- 多锁协调
- 优雅关闭

```bash
cd example/lock_advanced
go run main.go
```

## 分布式锁 API

### 创建锁

```go
import deimos "github.com/marsevilspirit/deimos-client"

client := deimos.NewClient([]string{"http://127.0.0.1:4001"})

// 使用默认选项
lock := client.NewDistributedLock("/locks/my-resource", "node-id")

// 使用自定义选项
lock := client.NewDistributedLock("/locks/my-resource", "node-id",
    deimos.WithLockTTL(30*time.Second),        // 锁的生存时间
    deimos.WithRenewalPeriod(10*time.Second),  // 续约间隔
    deimos.WithAutoRenewal(true))              // 是否自动续约
```

### 基本操作

```go
ctx := context.Background()

// 获取锁（阻塞直到成功）
err := lock.Lock(ctx)

// 尝试获取锁（非阻塞）
err := lock.TryLock(ctx)

// 带重试的获取锁
err := lock.LockWithRetry(ctx, 100*time.Millisecond, 10)

// 检查锁状态
isHeld := lock.IsHeld()

// 续约锁
err := lock.Renew(ctx)

// 释放锁
err := lock.Unlock(ctx)
```

### 便捷方法

```go
// 自动处理锁的获取和释放
err := lock.WithLock(ctx, func() error {
    // 在锁保护下执行的代码
    return doSomeWork()
})
```

### 自动续约

```go
// 启动自动续约
lock.StartAutoRenewal(ctx, 10*time.Second)
```

### 锁信息

```go
info := lock.Info()
fmt.Printf("Key: %s, Held: %t, TTL: %v\n", 
    info.Key, info.Held, info.TTL)
```

## 最佳实践

### 1. 选择合适的 TTL

- **短 TTL (5-30秒)**: 适用于快速操作，减少死锁风险
- **长 TTL (1-5分钟)**: 适用于长时间操作，配合自动续约

### 2. 使用自动续约

对于长时间操作，建议启用自动续约：

```go
opts := &deimos.LockOptions{
    TTL:           30 * time.Second,
    RenewalPeriod: 10 * time.Second,  // TTL 的 1/3
    AutoRenewal:   true,
}
```

### 3. 错误处理

```go
if err := lock.Lock(ctx); err != nil {
    if errors.Is(err, deimos.ErrLockNotAcquired) {
        // 锁被其他节点持有
    } else {
        // 其他错误（网络、服务器等）
    }
}
```

### 4. 优雅关闭

```go
// 使用 context 控制生命周期
ctx, cancel := context.WithCancel(context.Background())

// 在程序退出时取消 context
defer cancel()

// 确保释放锁
defer func() {
    if lock.IsHeld() {
        lock.Unlock(context.Background())
    }
}()
```

### 5. 避免死锁

- 总是设置合理的 TTL
- 使用 context 设置超时
- 按固定顺序获取多个锁
- 及时释放不需要的锁

## 故障排除

### 常见错误

1. **ErrLockNotAcquired**: 锁被其他节点持有
   - 检查是否有其他实例在运行
   - 增加重试次数或间隔

2. **ErrLockNotHeld**: 尝试操作未持有的锁
   - 检查锁是否已过期
   - 确保在正确的时机调用操作

3. **ErrLockExpired**: 锁已过期
   - 增加 TTL 时间
   - 启用自动续约
   - 检查网络延迟

### 调试技巧

1. 启用详细日志
2. 监控锁的状态变化
3. 检查 Deimos 服务器日志
4. 使用 curl 直接检查锁状态：

```bash
# 检查锁状态
curl http://127.0.0.1:4001/keys/locks/my-resource

# 查看所有锁
curl http://127.0.0.1:4001/keys/locks?recursive=true
```

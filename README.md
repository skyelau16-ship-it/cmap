# cmap

基于 `sync.RWMutex` 的并发安全泛型 map，支持任意 `comparable` 类型的 key 和任意类型的 value。

## 使用

```go
m := cmap.New[string, int]()

m.Set("a", 1)
v, ok := m.Get("a") // (1, true)
m.Delete("a")
n := m.Len() // 0
```

## 设计要点

- 读操作（`Get`/`Len`）持共享读锁，写操作（`Set`/`Delete`）持独占写锁；读多写少场景收益最大。
- `Get` 返回 `(V, bool)`，区分「key 不存在」和「零值」。
- 内嵌 `sync.RWMutex`，实例勿按值复制，始终用指针传递。

## 测试

```bash
go test -race ./...
```

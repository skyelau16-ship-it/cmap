// Package cmap 提供一个基于 sync.RWMutex 的并发安全泛型 map。
package cmap

import "sync"

// ConcurrentMap 是并发安全的 map。读操作（Get/Len）共享读锁，写操作（Set/Delete）持独占写锁。
//
// 注意：内部含有 sync.RWMutex，实例使用指针传递，不要按值复制。
type ConcurrentMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// New 创建一个空的 ConcurrentMap。
func New[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{data: make(map[K]V)}
}

// Get 返回 key 对应的值，以及该 key 是否存在。
func (m *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

// Set 写入 key/value，若 key 已存在则覆盖。
func (m *ConcurrentMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[K]V)
	}
	m.data[key] = value
}

// Delete 删除 key，key 不存在时是空操作。
func (m *ConcurrentMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// Len 返回当前 map 中的键数量。
func (m *ConcurrentMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

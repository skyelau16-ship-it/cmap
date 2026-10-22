package cmap

import (
	"sync"
	"testing"
)

func TestBasicOperations(t *testing.T) {
	m := New[string, int]()

	if got := m.Len(); got != 0 {
		t.Fatalf("初始 Len() = %d, want 0", got)
	}

	// Get 不存在的 key：返回零值和 false
	if v, ok := m.Get("missing"); ok || v != 0 {
		t.Fatalf("Get(missing) = (%d, %v), want (0, false)", v, ok)
	}

	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("a", 100) // 覆盖
	if got := m.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
	if v, ok := m.Get("a"); !ok || v != 100 {
		t.Fatalf("Get(a) = (%d, %v), want (100, true)", v, ok)
	}

	m.Delete("a")
	m.Delete("not-exist") // no-op
	if _, ok := m.Get("a"); ok {
		t.Fatal("删除后 Get(a) 仍返回 true")
	}
	if got := m.Len(); got != 1 {
		t.Fatalf("Len() = %d, want 1", got)
	}
}

func TestConcurrentAccess(t *testing.T) {
	m := New[int, string]()
	const goroutines = 50
	const iterations = 200

	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				key := id*1000 + i
				m.Set(key, "v")
				m.Get(key)
				m.Len()
				if i%2 == 0 {
					m.Delete(key)
				}
			}
		}(g)
	}
	wg.Wait()

	// 每个 goroutine 留下 iterations/2 个奇数下标的 key
	want := goroutines * (iterations / 2)
	if got := m.Len(); got != want {
		t.Fatalf("并发写入后 Len() = %d, want %d", got, want)
	}
}

package lrucache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// capacity is 1
func Test_1(t *testing.T) {
	ast := assert.New(t)

	cache := NewLRUCache(1)

	cache.Put(0, 0)
	// [(0,0)]

	cache.Put(1, 1)
	// [(1,1)]

	cache.Put(2, 2)
	// [(2,2)]
	ret, err := cache.Get(1)
	t.Log(ast.Equal(-1, ret, "get 1 from [(2,2)]"))
	t.Log(ast.Equal("Not Found", err.Error(), "get 1 from [(2,2)]"))
	// [(1,1), (2,2), (0,0)]
	t.Log(cache.DumpKeys())
}

// capacity is 0
func Test_2(t *testing.T) {
	ast := assert.New(t)

	cache := NewLRUCache(0)

	b, err := cache.Put(1, 1)
	t.Log(ast.Equal(false, b, "put 1 for empty cache"))
	t.Log(ast.Equal("Cache capcity is 0", err.Error(), "put 1 for empty cache"))

	ret, err := cache.Get(1)
	t.Log(ast.Equal(-1, ret, "get 1 for empty cache"))
	t.Log(ast.Equal("Cache is empty or capcity is 0", err.Error(), "put 1 for empty cache"))
}

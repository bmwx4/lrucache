# lrucache

### usage
```go
import "github.com/bmwx4/lrucache"
func main(){
    c := lrucache.NewLRUCache(2)
    c.Put(1)
    ...
    c.Get(1)
    lrucache.CleanLRUCache()
}
```

### test case

测试容量为1 的场景,添加和查询node:
```sh
go test -v -run Test_1
=== RUN   Test_1
--- PASS: Test_1 (0.00s)
    lru_cache_test.go:24: true
    lru_cache_test.go:25: true
    lru_cache_test.go:27: [{2 2}]
PASS
ok  	lru-cache	0.013s
```

测试容量为0 的场景,添加和查询node:
```sh
go test -v -run Test_2
=== RUN   Test_2
--- PASS: Test_2 (0.00s)
    lru_cache_test.go:37: true
    lru_cache_test.go:38: true
    lru_cache_test.go:41: true
    lru_cache_test.go:42: true
PASS
ok  	lru-cache	0.013s
```
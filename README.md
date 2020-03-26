[![Build Status](https://travis-ci.com/iamjinlei/goringbuffer.svg?branch=master)](https://travis-ci.com/iamjinlei/goringbuffer)

# Ring Buffer

A simple thread safe ring buffer implementation that supports add and iterate operation. The position in ring buffer does not necessarily indicate data freshness at the tiem when data are added concurrently, but eventually would. An example of such scenario:


Assume we have following ring buffer

```
+----+----+----+----+----+----
| p0 | p1 | p2 | p3 | p4 |
+----+----+----+----+----+----
```


And the code is


```
r := goringbuffer.New(10)
go r.Add(100) // to p0 with old value 10
go r.Add(200) // to p1 with old value 20
sum := []int{0}
go r.Do(func(e interface{}){
   sum[0] += e.(int)
})
```

Even though the above code may settle 100 into p0 and 200 into p1 eventually, sum[0] may be 30 (10 + 20), 300 (100 + 200), 210 (10 + 200) or 120 (100 + 20).

# Perf

Intel Core i5 2.3GHz

#### size = 16
```
BenchmarkAdd-8          50000000                36.1 ns/op
BenchmarkDo-8           20000000                104 ns/op
```

#### size = 64
```
BenchmarkAdd-8          30000000                36.7 ns/op
BenchmarkDo-8            5000000                322 ns/op
```

package buffer

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	r := New(10)
	r.Add(1)
	r.Add(2)
	r.Add(3)
	sum := 0
	r.Do(func(e interface{}) {
		sum += e.(int)
	})
	assert.Equal(t, 6, sum)
	assert.Equal(t, int32(10), r.Capacity())

	r.Add(4)
	r.Add(5)
	r.Add(6)
	sum = 0
	r.Do(func(e interface{}) {
		sum += e.(int)
	})
	assert.Equal(t, 21, sum)

	r.Add(7)
	r.Add(8)
	r.Add(9)
	r.Add(10)
	r.Add(11)
	r.Add(12)
	sum = 0
	r.Do(func(e interface{}) {
		sum += e.(int)
	})
	assert.Equal(t, 75, sum)
}

func TestStressRun(t *testing.T) {
	add := func(r *Ring, d time.Duration, wg *sync.WaitGroup) {
		t.Log("add thread started")
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				wg.Done()
				return
			default:
				r.Add(rand.Int())
			}
		}
	}

	sum := func(r *Ring, d time.Duration, wg *sync.WaitGroup) {
		t.Log("sum thread started")
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				wg.Done()
				return
			default:
				data := 0
				r.Do(func(e interface{}) {
					data += e.(int)
				})
			}
		}
	}

	d := 30 * time.Second
	n := 4

	r := New(10)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go add(r, d, &wg)
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go sum(r, d, &wg)
	}
	wg.Wait()
}

const size = 64

func BenchmarkAdd(b *testing.B) {
	r := New(size)
	for i := 0; i < b.N; i++ {
		r.Add(i)
	}
}

func BenchmarkDo(b *testing.B) {
	r := New(size)
	for i := 0; i < size; i++ {
		r.Add(i)
	}
	for i := 0; i < b.N; i++ {
		data := 0
		r.Do(func(e interface{}) {
			data += e.(int)
		})
	}
}

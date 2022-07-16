package lru

import (
	"fmt"
	"sync"
	"testing"

	"gotest.tools/assert"
)

type testEntry int32

func (e testEntry) Key() string {
	return fmt.Sprintf("%d", e)
}

func (e testEntry) Size() int64 {
	return 4
}

func TestCache(t *testing.T) {
	c := NewCache(32)
	debug := make(chan struct{})
	c.debug = debug

	c.Add(testEntry(1))
	c.Add(testEntry(2))
	c.Add(testEntry(3))
	c.Add(testEntry(4))
	c.Add(testEntry(5))
	c.Add(testEntry(6))
	c.Add(testEntry(7))
	c.Add(testEntry(8))

	assert.Equal(t, 8, len(c.entries))
	assert.Equal(t, int64(32), c.Size())

	e, ok := c.Get("1")
	assert.Equal(t, true, ok)
	assert.Equal(t, testEntry(1), e.(testEntry))
	<-c.debug // Updates history async.

	e, ok = c.Get("10")
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, e)

	c.Add(testEntry(9))
	<-c.debug // Resizes async.

	e, ok = c.Get("2")
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, e)

	assert.Equal(t, int64(32), c.Size())
}

func BenchmarkAddConcurrent(b *testing.B) {
	c := NewCache(30)
	debug := make(chan struct{})
	c.debug = debug

	var wg sync.WaitGroup

	for n := 0; n < b.N; n++ {
		c.reset()

		for g := 0; g < 8; g++ {
			wg.Add(1)
			go func(id int) {
				c.Add(testEntry(id))
				wg.Done()
			}(g)
		}

		wg.Wait()
		<-debug // Resizes async.
	}

	assert.Equal(b, int64(28), c.Size())
}

func BenchmarkAdd(b *testing.B) {
	c := NewCache(4)

	for n := 0; n < b.N; n++ {
		c.Add(testEntry(10))
	}

	assert.Equal(b, int64(4), c.Size())
}

func BenchmarkResize(b *testing.B) {
	c := NewCache(5)
	debug := make(chan struct{})
	c.debug = debug

	c.Add(testEntry(-1))

	assert.Equal(b, int64(4), c.Size())

	for n := 0; n < b.N; n++ {
		c.Add(testEntry(n))
		<-debug
	}

	assert.Equal(b, int64(4), c.Size())
}

package lru

import (
	"sync"
	"sync/atomic"

	"github.com/rchilly/lru/internal"
)

type Entry interface {
	Key() string
	Size() int64
}

type entry struct {
	Entry

	recurrer internal.Recurrer
}

// TODO: Make the lock/unlock calls no-ops if user
// doesn't want thread safety. Compare benchmarks.

type Cache struct {
	entries map[string]entry
	emu     *sync.RWMutex

	history *internal.History
	hmu     *sync.Mutex

	resizing chan struct{}
	maxSize  int64
	size     *int64

	debug chan struct{}
}

func NewCache(maxSize int64) *Cache {
	if maxSize <= 0 {
		panic("max size must be greater than 0")
	}

	cache := &Cache{
		entries: make(map[string]entry),
		emu:     &sync.RWMutex{},

		history: &internal.History{},
		hmu:     &sync.Mutex{},

		resizing: make(chan struct{}, 1),
		maxSize:  maxSize,
		size:     new(int64),
	}

	return cache
}

// TODO: Get and Add become private, so I can
// inject a monitor/spy of some kind to track
// whether cache is all done.

func (c *Cache) Get(key string) (Entry, bool) {
	c.emu.RLock()
	e, ok := c.entries[key]
	c.emu.RUnlock()

	if !ok {
		return nil, false
	}

	go c.recur(e)

	return e.Entry, true
}

func (c *Cache) recur(e entry) {
	c.hmu.Lock()
	e.recurrer.Recur()
	c.hmu.Unlock()

	if c.debug != nil {
		c.debug <- struct{}{}
	}
}

func (c *Cache) Add(e Entry) {
	if e == nil {
		return
	}

	key := e.Key()

	c.emu.RLock()
	old, overwriting := c.entries[key]
	c.emu.RUnlock()

	if overwriting {
		old.Entry = e

		c.emu.Lock()
		c.entries[key] = old
		c.emu.Unlock()

		go c.recur(old)

		return
	}

	size := atomic.AddInt64(c.size, e.Size())
	if size > c.maxSize {
		go c.resize()
	}

	c.hmu.Lock()
	recurrer := c.history.Track(key)
	c.hmu.Unlock()

	c.emu.Lock()
	c.entries[key] = entry{
		Entry:    e,
		recurrer: recurrer,
	}
	c.emu.Unlock()
}

func (c *Cache) resize() {
	select {
	case c.resizing <- struct{}{}:
	default:
		return
	}

	trim := atomic.LoadInt64(c.size) - c.maxSize
	trimmed := trim

	var keys []string

	c.hmu.Lock()
	for trimmed > 0 && !c.history.Empty() {
		key := c.history.Forget()

		c.emu.RLock()
		e, ok := c.entries[key]
		c.emu.RUnlock()

		if !ok {
			continue
		}

		trimmed -= e.Size()
		keys = append(keys, e.Key())
	}
	c.hmu.Unlock()

	c.emu.Lock()
	for _, k := range keys {
		delete(c.entries, k)
	}
	c.emu.Unlock()

	atomic.AddInt64(c.size, -(trim - trimmed))

	<-c.resizing

	if c.debug != nil {
		c.debug <- struct{}{}
	}
}

func (c Cache) Size() int64 {
	return atomic.LoadInt64(c.size)
}

func (c *Cache) reset() {
	for k := range c.entries {
		delete(c.entries, k)
	}

	c.history = &internal.History{}

	c.size = new(int64)
}

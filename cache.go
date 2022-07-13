package lru

import (
	"github.com/rchilly/lru/internal"
)

type Entry interface {
	Key() string
	Size() uint64
}

type entry struct {
	Entry

	seen *internal.Seen
}

type Cache struct {
	entries map[string]entry
	head    *internal.Seen
	tail    *internal.Seen

	size, maxSize uint64
}

func NewCache(maxEntries, maxSize uint64) *Cache {
	return &Cache{
		entries: make(map[string]entry, maxSize),
	}
}

func (c *Cache) Get(key string) (Entry, bool) {
	e, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	c.head = c.head.SawAgain(e.seen)

	return e, true
}

func (c *Cache) Add(e Entry) {
	if e == nil {
		return
	}

	c.size += e.Size()

	for c.size > c.maxSize && c.tail != nil {
		tail := c.tail
		c.tail = tail.After()

		rm, ok := c.entries[tail.Key()]
		if !ok {
			continue
		}

		c.size -= rm.Size()
	}

	seen := c.head.Saw(e.Key())

	c.entries[e.Key()] = entry{
		Entry: e,
		seen:  seen,
	}

	c.head = seen
}

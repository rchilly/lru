package lru

import (
	"fmt"

	"github.com/rchilly/lru/internal"
)

type Entry interface {
	Key() string
	Size() uint64
	Value() interface{}
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

func (c Cache) String() string {
	return fmt.Sprintf(
		"\nsize: %d\nhead: %s\ntail: %s\n",
		c.size,
		c.head.Key(),
		c.tail.Key(),
	)
}

func NewCache(maxSize uint64) *Cache {
	return &Cache{
		entries: make(map[string]entry),
		maxSize: maxSize,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	e, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	c.head = c.head.SawAgain(e.seen)

	if c.tail == e.seen {
		c.tail = e.seen.After()
	}

	return e.Value(), true
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
		delete(c.entries, tail.Key())
	}

	seen := c.head.Saw(e.Key())

	c.entries[e.Key()] = entry{
		Entry: e,
		seen:  seen,
	}

	c.head = seen

	if c.tail == nil {
		c.tail = seen
	}
}

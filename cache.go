package lru

import (
	"fmt"

	"github.com/rchilly/lru/internal"
)

type Entry interface {
	Key() string
	Size() uint64
}

type entry struct {
	Entry

	recurrer internal.Recurrer
}

type Cache struct {
	entries map[string]entry
	history *internal.History

	size, maxSize uint64
}

func NewCache(maxSize uint64) *Cache {
	return &Cache{
		entries: make(map[string]entry),
		history: &internal.History{},
		maxSize: maxSize,
	}
}

func (c *Cache) Get(key string) (Entry, bool) {
	if c == nil {
		return nil, false
	}

	e, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	e.recurrer.Recur()

	return e.Entry, true
}

func (c *Cache) Add(e Entry) {
	if e == nil {
		return
	}

	c.size += e.Size()

	for c.size > c.maxSize && !c.history.Empty() {
		key := c.history.Forget()
		rm, ok := c.entries[key]
		if !ok {
			continue
		}

		c.size -= rm.Size()
		delete(c.entries, key)
	}

	recurrer := c.history.Track(e.Key())

	c.entries[e.Key()] = entry{
		Entry:    e,
		recurrer: recurrer,
	}
}

func (c Cache) String() string {
	mru, lru := c.history.Bookends()

	return fmt.Sprintf(
		"\nsize: %d\nmru: %s\nlru: %s\n",
		c.size,
		mru,
		lru,
	)
}

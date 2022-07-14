package lru

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

type testEntry int32

func (e testEntry) Key() string {
	return fmt.Sprintf("%d", e)
}

func (e testEntry) Size() uint64 {
	return 4
}

func TestCache(t *testing.T) {
	c := NewCache(32)
	c.Add(testEntry(1))
	c.Add(testEntry(2))
	c.Add(testEntry(3))
	c.Add(testEntry(4))
	c.Add(testEntry(5))
	c.Add(testEntry(6))
	c.Add(testEntry(7))
	c.Add(testEntry(8))

	assert.Equal(t, 8, len(c.entries))
	t.Logf("after 1-8 added: %v", c.history.LatestN(10))

	e, ok := c.Get("1")
	assert.Equal(t, true, ok)
	assert.Equal(t, testEntry(1), e.(testEntry))
	t.Logf("after 1 requested: %s", c.history.LatestN(10))

	e, ok = c.Get("10")
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, e)

	c.Add(testEntry(9))
	t.Logf("after 9 added: %s", c.history.LatestN(10))

	e, ok = c.Get("2")
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, e)
}

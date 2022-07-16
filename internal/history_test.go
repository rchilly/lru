package internal

import (
	"testing"

	"gotest.tools/assert"
)

func TestTrack(t *testing.T) {
	h := &History{}

	a := h.Track("A")
	b := h.Track("B")
	_ = h.Track("C")
	_ = h.Track("D")
	e := h.Track("E")

	assert.DeepEqual(t, []string{"E", "D", "C", "B", "A"}, h.LatestN(5))

	b.Recur()

	assert.DeepEqual(t, []string{"B", "E", "D", "C", "A"}, h.LatestN(5))

	b.Recur()

	assert.DeepEqual(t, []string{"B", "E", "D", "C", "A"}, h.LatestN(5))

	a.Recur()

	assert.DeepEqual(t, []string{"A", "B", "E", "D", "C"}, h.LatestN(5))

	b.Recur()

	assert.DeepEqual(t, []string{"B", "A", "E", "D", "C"}, h.LatestN(5))

	a2 := h.Track("A")
	_ = h.Track("B")

	assert.DeepEqual(t, []string{"B", "A", "B", "A", "E", "D", "C"}, h.LatestN(10))

	a2.Recur()

	assert.DeepEqual(t, []string{"A", "B", "B", "A", "E", "D", "C"}, h.LatestN(10))

	e.Recur()

	assert.DeepEqual(t, []string{"E", "A", "B", "B", "A", "D", "C"}, h.LatestN(10))
}

func TestForget(t *testing.T) {
	h := &History{}
	one := h.Track("1")
	_ = h.Track("2")
	_ = h.Track("3")

	one.Recur()

	forgotten := h.Forget()

	assert.Equal(t, "2", forgotten)
	assert.DeepEqual(t, []string{"1", "3"}, h.LatestN(5))
}

func BenchmarkRecur(tb *testing.B) {
	h := &History{}
	a := h.Track("A")
	b := h.Track("B")
	c := h.Track("C")

	// C B A
	// B C A
	// ...etc

	tb.ResetTimer()

	for n := 0; n < tb.N; n++ {
		b.Recur()
		c.Recur()
	}

	assert.Equal(tb, c.(recurrer).l, h.head)
	assert.Equal(tb, a.(recurrer).l, h.tail)
}

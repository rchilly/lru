package internal

import (
	"testing"

	"gotest.tools/assert"
)

func TestSeen(t *testing.T) {
	var s *Seen
	a := s.Saw("A")
	b := a.Saw("B")
	s = b.Saw("C")
	s = s.Saw("D")
	s = s.Saw("E")

	assert.DeepEqual(t, []string{"E", "D", "C", "B", "A"}, s.lastN(5))

	s = s.SawAgain(b)

	assert.DeepEqual(t, []string{"B", "E", "D", "C", "A"}, s.lastN(5))

	s = s.SawAgain(s)

	assert.DeepEqual(t, []string{"B", "E", "D", "C", "A"}, s.lastN(5))

	a = s.SawAgain(a)

	assert.DeepEqual(t, []string{"A", "B", "E", "D", "C"}, a.lastN(5))

	s = nil
	s = s.SawAgain(a)

	assert.DeepEqual(t, []string{"A"}, s.lastN(5))
}

package internal

// History tracks occurrences and recurrences as an ordered
// sequence of objects. A zero value is ready-to-use.
//
// An occurrence joins the front of history as a new object
// while a recurrence is just the movement of an old object
// back to the front.
type History struct {
	head *ll
	tail *ll
}

type ll struct {
	s     string
	left  *ll
	right *ll
}

// Recurrer repeats a former occurrence, so that
// it returns to the front of its parent history.
type Recurrer interface {
	Recur()
}

// Stores the address of an occurrence in its parent
// history so that both can be updated to execute its
// recurrence.
type recurrer struct {
	l *ll
	h *History
}

func (r recurrer) Recur() {
	l, h := r.l, r.h

	// 1. Check for no-op scenarios.

	if h == nil || h.head == l || l == nil {
		return
	}

	// 2. Pluck l from its current context:
	// - If l has a left, point it to l's right.
	// - If l has a right, point it to l's left.
	// - If l is the parent history's tail, make
	//   its left the new tail.

	left, right := l.left, l.right
	if left != nil {
		left.right = right
	}
	if right != nil {
		right.left = left
	}
	if h.tail == l {
		h.tail = l.left
	}

	// 3. Prepare l for its new context as head

	l.left = nil
	l.right = h.head

	// 4. Add l as new head to the left of old head

	if h.head != nil {
		h.head.left = l
	}

	h.head = l
}

// Track pushes a new occurrence to the front of history
// and returns an interface to promote that same occurrence
// back to the front of history as a recurrence.
func (h *History) Track(occurrence string) Recurrer {
	newbie := &ll{
		s:     occurrence,
		right: h.head,
	}

	if h.head == nil {
		h.tail = newbie
	} else {
		h.head.left = newbie
	}

	h.head = newbie

	return recurrer{newbie, h}
}

// Forget pops and returns the oldest occurrence from history.
func (h *History) Forget() (oldest string) {
	tail := h.tail
	if tail == nil {
		return ""
	}

	left := tail.left

	if left == nil {
		h.head = nil
	} else {
		left.right = nil
	}

	h.tail = left

	return tail.s
}

// LatestN returns a slice of the n latest occurrences or
// recurrences in history, in order, where index 0 is the
// most recent.
func (h *History) LatestN(n int) (occurrences []string) {
	occurrences = make([]string, 0, n)
	l := h.head
	for l != nil && n > 0 {
		occurrences = append(occurrences, l.s)
		l = l.right
		n--
	}

	return
}

// Bookends peeks at the latest and oldest occurrences
// in history.
func (h *History) Bookends() (latest, oldest string) {
	if h.head != nil {
		latest = h.head.s
	}

	if h.tail != nil {
		oldest = h.tail.s
	}

	return
}

func (h *History) Empty() bool {
	return h.head == nil
}

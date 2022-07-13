package internal

type Seen struct {
	key    string
	after  *Seen
	before *Seen
}

func (s *Seen) After() *Seen {
	return s.after
}

func (s *Seen) Key() string {
	return s.key
}

func (s *Seen) Saw(key string) *Seen {
	newbie := &Seen{
		key:    key,
		before: s,
	}

	if s != nil {
		s.after = newbie
	}

	return newbie
}

func (s *Seen) SawAgain(already *Seen) *Seen {
	if s == already || already == nil {
		return s
	}

	after, before := already.after, already.before
	if after != nil {
		after.before = before
	}
	if before != nil {
		before.after = after
	}

	already.before = s

	return already
}

func (s *Seen) lastN(n int) []string {
	got := make([]string, 0, n)
	for s != nil {
		got = append(got, s.key)
		s = s.before
	}

	return got
}

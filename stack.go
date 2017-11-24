package shortlivedpool

import (
	"sync"
	"time"
)

const (
	defaultMinSize = 64
	maxTimeInStack = 60 * time.Second
)

type item struct {
	x  interface{}
	ts time.Time
}

type Stack struct {
	sync.Mutex
	vec []*item
}

func (s *Stack) Put(x interface{}) {
	now := time.Now()
	s.Lock()
	s.vec = append(s.vec, &item{
		x:  x,
		ts: now,
	})
	if len(s.vec) == 1 {
		s.Unlock()
		return
	}
	if s.vec[0].ts.Add(maxTimeInStack).Before(now) {
		// Take out the oldest element
		s.vec[0] = nil
		s.vec = s.vec[1:]
		if len(s.vec) > defaultMinSize {
			s.shrink()
		}
	}
	s.Unlock()
}

func (s *Stack) Pop() interface{} {
	s.Lock()
	l := len(s.vec)
	if l == 0 {
		s.Unlock()
		return nil
	}
	x := s.vec[l-1].x
	s.vec[l-1] = nil
	s.vec = s.vec[:l-1]
	s.Unlock()
	return x
}

func (s *Stack) shrink() {
	if cap(s.vec) > len(s.vec)*4 {
		newVec := make([]*item, len(s.vec), len(s.vec)*2)
		copy(newVec, s.vec)
		s.vec = newVec
	}
}

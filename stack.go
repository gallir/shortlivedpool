package shortlivedpool

import (
	"sync"
	"time"
)

const (
	defaultMinSize             = 64
	maxSecondsInStack          = 60
	minSecondsBetweenEvictions = 1
)

type item struct {
	x  interface{}
	ts int64
}

type Stack struct {
	sync.Mutex
	vec          []item
	nextEviction int64
}

func (s *Stack) Put(x interface{}) {
	now := time.Now().Unix()
	s.Lock()
	s.vec = append(s.vec, item{
		x:  x,
		ts: now,
	})
	if len(s.vec) == 1 {
		s.Unlock()
		return
	}
	if s.nextEviction < now && s.vec[0].ts+maxSecondsInStack < now {
		// Evict the oldest element
		s.vec[0].x = nil
		s.vec = s.vec[1:]
		s.nextEviction = now + minSecondsBetweenEvictions
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
	s.vec[l-1].x = nil
	s.vec = s.vec[:l-1]
	s.Unlock()
	return x
}

func (s *Stack) shrink() {
	if cap(s.vec) > len(s.vec)*4 {
		newVec := make([]item, len(s.vec), len(s.vec)*2)
		copy(newVec, s.vec)
		s.vec = newVec
	}
}

package shortlivedpool

import (
	"sync"
	"time"
)

const (
	defaultMinSize             = 64
	maxSecondsInStack          = 60
	minSecondsBetweenEvictions = 15
)

type item struct {
	x  interface{}
	ts int64
}

// EvictionStack remove elements that are older
// than maxSecondsInStack seconds in the stack
type EvictionStack struct {
	sync.Mutex
	vec          []item
	nextEviction int64 // To avoid checking for evictions for every put
}

// Put pushes en element into the stack
func (s *EvictionStack) Put(x interface{}) {
	now := time.Now().Unix()
	s.Lock()
	s.vec = append(s.vec, item{
		x:  x,
		ts: now,
	})

	if len(s.vec) == 1 || s.nextEviction > now {
		s.Unlock()
		return
	}

	// Evict the oldest elements
	for s.vec[0].ts+maxSecondsInStack < now {
		s.vec[0].x = nil
		s.vec = s.vec[1:]
	}

	if len(s.vec) > defaultMinSize {
		s.shrink()
	}

	s.nextEviction = now + minSecondsBetweenEvictions
	s.Unlock()
}

// Pop gets the last element inserted
func (s *EvictionStack) Pop() interface{} {
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

func (s *EvictionStack) shrink() {
	if cap(s.vec) > len(s.vec)*4 {
		newVec := make([]item, len(s.vec), len(s.vec)*2)
		copy(newVec, s.vec)
		s.vec = newVec
	}
}

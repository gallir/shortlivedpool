package shortlivedpool

import "time"

const (
	defaultMinSize = 64
	maxTimeInStack = 60 * time.Second
)

type item struct {
	x  interface{}
	ts time.Time
}

type Stack struct {
	oldest time.Time
	vec    []item
}

var ()

func (s Stack) Peek() interface{} {
	return s.vec[len(s.vec)-1].x
}

func (s Stack) Len() int {
	return len(s.vec)
}

func (s *Stack) Put(x interface{}) {
	now := time.Now()
	s.vec = append(s.vec, item{
		x:  x,
		ts: now,
	})
	if len(s.vec) == 1 {
		s.oldest = now
		return
	}
	if s.oldest.Add(maxTimeInStack).Before(now) {
		s.shrink()
	}
}

func (s *Stack) Pop() interface{} {
	l := len(s.vec)
	if l == 0 {
		return nil
	}
	x := s.vec[l-1].x
	s.vec = s.vec[:l-1]
	return x
}

func (s *Stack) shrink() {
	pos := 0
	now := time.Now()
	base := now.Add(-maxTimeInStack)
	for i := range s.vec {
		if s.vec[i].ts.Before(base) {
			pos++
			s.vec[i].x = nil
		} else {
			break
		}
	}

	if pos == 0 {
		return
	}

	s.vec = s.vec[pos:]
	if len(s.vec) > 0 {
		s.oldest = s.vec[0].ts
	} else {
		s.oldest = now
	}

	if cap(s.vec) > len(s.vec)*4 {
		newVec := make([]item, len(s.vec), len(s.vec)*2)
		copy(newVec, s.vec)
		s.vec = newVec
	}
}

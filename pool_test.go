package shortlivedpool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var (
	aString = "A string"
)

func Test_PutGet(t *testing.T) {
	p := &Pool{
		New: func() interface{} {
			return aString
		},
	}

	got := p.Get().(string)
	if got != aString {
		t.Errorf("Pool.Get() = %s, want %s", got, aString)
	}

	another := "Another string"
	p.Put(another)
	p.Put(got)

	got = p.Get().(string)
	if got != aString {
		t.Errorf("Pool.Get() = %s, want %s", got, aString)
	}

	got = p.Get().(string)
	if got != another {
		t.Errorf("Pool.Get() = %s, want %s", got, another)
	}
}

type pooler interface {
	Put(x interface{})
	Get() interface{}
}

func Test_Benchmark(t *testing.T) {
	p1 := &sync.Pool{
		New: func() interface{} {
			return aString
		},
	}

	p2 := &Pool{
		New: func() interface{} {
			return aString
		},
	}

	tests := []pooler{p1, p2}
	t.Log("Running benchmarks")
	for i, p := range tests {
		start := time.Now()
		for i := 0; i < 1000; i++ {
			p.Put(fmt.Sprintf("N %d", i))
		}
		for i := 0; i < 10000000; i++ {
			s := p.Get().(string)
			p.Put(s)
		}
		fmt.Println("time", i, time.Since(start))
	}
}

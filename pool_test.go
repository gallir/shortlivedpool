package shortlivedpool

import (
	"sync"
	"testing"
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

func simple(b *testing.B, p pooler) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Put(1)
			p.Get()
		}
	})
}

func overflow(b *testing.B, p pooler) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for b := 0; b < 100; b++ {
				p.Put(1)
			}
			for b := 0; b < 100; b++ {
				p.Get()
			}
		}
	})
}

func BenchmarkSyncPoolSimple(b *testing.B) {
	p := &sync.Pool{}

	simple(b, p)
}

func BenchmarkShortLivedPoolSimple(b *testing.B) {
	p := &Pool{}

	simple(b, p)
}

func BenchmarkSyncPoolOverflow(b *testing.B) {
	p := &sync.Pool{}

	overflow(b, p)
}

func BenchmarkShortLivedPoolOverflow(b *testing.B) {
	p := &Pool{}

	overflow(b, p)
}

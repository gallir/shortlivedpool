package shortlivedpool

type noCopy struct{}

type Pool struct {
	noCopy noCopy
	stack  Stack
	New    func() interface{}
}

// Put adds x to the pool.
func (p *Pool) Put(x interface{}) {
	p.stack.Put(x)
}

// Get selects the most recent used item from the Pool,
// removes it from the Pool, and returns it to the caller.
// Get may choose to ignore the pool and treat it as empty.
// Callers should not assume any relation between values passed to Put and
// the values returned by Get.
//
// If Get would otherwise return nil and p.New is non-nil, Get returns
// the result of calling p.New.
func (p *Pool) Get() (x interface{}) {
	x = p.stack.Pop()
	if x == nil && p.New != nil {
		return p.New()
	}
	return
}

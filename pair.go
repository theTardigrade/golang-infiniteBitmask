package infiniteBitmask

type Pair struct {
	inner       pairInner
	innerInited bool
}

type pairInner struct {
	name  string
	value *Value
}

func newPair(name string, value *Value) *Pair {
	return &Pair{
		inner: pairInner{
			name:  name,
			value: value,
		},
		innerInited: true,
	}
}

func (p *Pair) checkInnerInited() {
	if p == nil {
		panic(ErrPointerNil)
	}

	if !p.innerInited {
		panic(ErrPairUninitialized)
	}
}

func (p *Pair) Name() string {
	p.checkInnerInited()

	return p.inner.name
}

func (p *Pair) Value() *Value {
	p.checkInnerInited()

	return p.inner.value
}

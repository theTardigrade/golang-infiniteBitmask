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

	return p.inner.value.Clone()
}

func (p *Pair) Clone() (p2 *Pair) {
	p.checkInnerInited()

	p2 = newPair(p.inner.name, p.inner.value.Clone())

	return
}

func (p *Pair) Equal(p2 *Pair) (result bool) {
	p.checkInnerInited()
	p2.checkInnerInited()

	if p.inner.name != p2.inner.name {
		return
	}

	if !p.inner.value.equal(p2.inner.value, false) {
		return
	}

	result = true

	return
}

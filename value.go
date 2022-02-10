package infiniteBitmask

import (
	"math/big"
	"sync"
)

type Value struct {
	inner       valueInner
	innerInited bool
	mutex       sync.RWMutex
}

type valueInner struct {
	number    *big.Int
	generator *Generator
}

func newValue(number *big.Int, generator *Generator) (v *Value) {
	v = &Value{}

	v.initInner(number, generator)

	return
}

func (v *Value) initInner(number *big.Int, generator *Generator) {
	if v.inner.number = new(big.Int); number != nil {
		v.inner.number.Set(number)
	}
	v.inner.generator = generator
	v.innerInited = true
}

func (v *Value) read(handler func()) {
	if v == nil {
		panic(ErrPointerNil)
	}

	var shouldWrite bool

	func() {
		defer v.mutex.RUnlock()
		v.mutex.RLock()

		if !v.innerInited {
			shouldWrite = true
			return
		}

		handler()
	}()

	if shouldWrite {
		v.write(handler)
	}
}

func (v *Value) write(handler func()) {
	if v == nil {
		panic(ErrPointerNil)
	}

	defer v.mutex.Unlock()
	v.mutex.Lock()

	if !v.innerInited {
		v.initInner(nil, nil)
	}

	handler()
}

func (v *Value) checkGeneratorMatch(v2 *Value) {
	if v.inner.generator != v2.inner.generator {
		panic(ErrValuesMismatched)
	}
}

func (v *Value) Combine(vs ...*Value) {
	v.write(func() {
		for _, v2 := range vs {
			v2.read(func() {
				v.checkGeneratorMatch(v2)

				v.inner.number.Or(v.inner.number, v2.inner.number)
			})
		}
	})
}

func (v *Value) Uncombine(vs ...*Value) {
	v.write(func() {
		for _, v2 := range vs {
			v2.read(func() {
				v.checkGeneratorMatch(v2)

				mask := new(big.Int)
				mask.Not(v2.inner.number)

				v.inner.number.And(v.inner.number, mask)
			})
		}
	})
}

func (v *Value) Contains(vs ...*Value) (result bool) {
	v.read(func() {
		intersection := new(big.Int)
		intersection.Set(v.inner.number)

		for _, v2 := range vs {
			v2.read(func() {
				v.checkGeneratorMatch(v2)

				intersection.And(intersection, v2.inner.number)
			})
		}

		result = intersection.Cmp(bigZero) == 1
	})

	return
}

func (v *Value) Clear() {
	v.write(func() {
		v.inner.number.Set(bigZero)
	})
}

func (v *Value) IsEmpty() (result bool) {
	v.read(func() {
		result = v.inner.number.Cmp(bigZero) == 0
	})

	return
}

func (v *Value) Clone() (v2 *Value) {
	v.read(func() {
		v2 = &Value{}

		v2.initInner(v.inner.number, v.inner.generator)
	})

	return
}

func (v *Value) equal(v2 *Value, checkGeneratorMatch bool) (result bool) {
	v.read(func() {
		v2.read(func() {
			if checkGeneratorMatch {
				v.checkGeneratorMatch(v2)
			}

			result = v.inner.number.Cmp(v2.inner.number) == 0
		})
	})

	return
}

func (v *Value) Equal(v2 *Value) (result bool) {
	result = v.equal(v2, true)

	return
}

func (v *Value) String() (result string) {
	v.read(func() {
		result = v.inner.number.Text(2)
	})

	return
}

func (v *Value) Number() (number *big.Int) {
	v.read(func() {
		number = v.inner.number
	})

	return
}

func (v *Value) Generator() (generator *Generator) {
	v.read(func() {
		generator = v.inner.generator
	})

	return
}

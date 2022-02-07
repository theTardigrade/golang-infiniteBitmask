package infiniteBitmask

import (
	"math/big"
	"sync"
)

type Value struct {
	number    *big.Int
	generator *Generator
	mutex     sync.RWMutex
}

func newValue(number int64, generator *Generator) (v *Value) {
	v = &Value{
		number:    big.NewInt(number),
		generator: generator,
	}

	return
}

func (v *Value) read(handler func(*Value)) {
	if v == nil {
		return
	}

	defer v.mutex.RUnlock()
	v.mutex.RLock()

	handler(v)
}

func (v *Value) write(handler func(*Value)) {
	if v == nil {
		return
	}

	defer v.mutex.Unlock()
	v.mutex.Lock()

	handler(v)
}

func (v *Value) Combine(vs ...*Value) {
	v.write(func(v *Value) {
		for _, v2 := range vs {
			v2.read(func(v2 *Value) {
				if v2.generator != v.generator {
					panic(ErrValuesMismatched)
				}

				v.number.Or(v.number, v2.number)
			})
		}
	})
}

func (v *Value) Uncombine(vs ...*Value) {
	v.write(func(v *Value) {
		for _, v2 := range vs {
			v2.read(func(v2 *Value) {
				if v2.generator != v.generator {
					panic(ErrValuesMismatched)
				}

				mask := new(big.Int)
				mask.Not(v2.number)

				v.number.And(v.number, mask)
			})
		}
	})
}

func (v *Value) Contains(vs ...*Value) (result bool) {
	v.read(func(v *Value) {
		intersection := new(big.Int)
		intersection.Set(v.number)

		for _, v2 := range vs {
			v2.read(func(v2 *Value) {
				if v2.generator != v.generator {
					panic(ErrValuesMismatched)
				}

				intersection.And(intersection, v2.number)
			})
		}

		result = intersection.Cmp(bigZero) == 1
	})

	return
}

func (v *Value) Clear() {
	v.write(func(v *Value) {
		v.number.Set(bigZero)
	})
}

func (v *Value) IsNotEmpty() (result bool) {
	v.read(func(v *Value) {
		result = v.number.Cmp(bigZero) != 0
	})

	return
}

func (v *Value) IsEmpty() (result bool) {
	result = !v.IsNotEmpty()

	return
}

func (v *Value) Clone() (v2 *Value) {
	v.read(func(v *Value) {
		n2 := new(big.Int)
		n2.Set(v.number)

		v2 = &Value{
			number:    n2,
			generator: v.generator,
		}
	})

	return
}

func (v *Value) Number() (number *big.Int) {
	v.read(func(v *Value) {
		number = v.number
	})

	return
}

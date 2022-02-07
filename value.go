package infiniteBitmask

import "math/big"

type Value struct {
	i *big.Int
}

func (v Value) Combine(vs ...Value) {
	for _, v2 := range vs {
		v.i.Or(v.i, v2.i)
	}
}

func (v Value) Uncombine(vs ...Value) {
	mask := new(big.Int)

	for _, v2 := range vs {
		mask.Not(v2.i)

		v.i.And(v.i, mask)
	}
}

func (v Value) Contains(v2 Value) bool {
	intersection := new(big.Int)
	intersection.And(v.i, v2.i)

	return intersection.Cmp(bigZero) == 1
}

func (v Value) Clear() {
	v.i.Set(bigZero)
}

func (v Value) IsEmpty() bool {
	return v.i.Cmp(bigZero) == 0
}

func (v Value) Clone() (v2 Value) {
	i2 := new(big.Int)
	i2.Set(v.i)

	v2.i = i2

	return
}

func (v Value) BigInt() *big.Int {
	return v.i
}

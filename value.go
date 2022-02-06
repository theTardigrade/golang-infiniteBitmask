package infiniteBitmask

import "math/big"

type Value struct {
	i *big.Int
}

func (v Value) Combine(v2 Value) {
	v.i.Or(v.i, v2.i)
}

func (v Value) Uncombine(v2 Value) {
	mask := new(big.Int)
	mask = v2.i.Not(v2.i)

	v.i.And(v.i, mask)
}

func (v Value) Contains(v2 Value) bool {
	intersection := new(big.Int)
	intersection.And(v.i, v2.i)

	return intersection.Cmp(bigZero) == 1
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

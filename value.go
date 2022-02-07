package infiniteBitmask

import "math/big"

type Value struct {
	i *big.Int
	g *Generator
}

func newValue(n int64, g *Generator) (v Value) {
	v.i = big.NewInt(n)

	return
}

func (v Value) Combine(vs ...Value) {
	for _, v2 := range vs {
		if v2.g != v.g {
			panic(ErrValuesMismatched)
		}

		v.i.Or(v.i, v2.i)
	}
}

func (v Value) Uncombine(vs ...Value) {
	mask := new(big.Int)

	for _, v2 := range vs {
		if v2.g != v.g {
			panic(ErrValuesMismatched)
		}

		mask.Not(v2.i)

		v.i.And(v.i, mask)
	}
}

func (v Value) Contains(vs ...Value) bool {
	intersection := new(big.Int)
	intersection.Set(v.i)

	for _, v2 := range vs {
		if v2.g != v.g {
			panic(ErrValuesMismatched)
		}

		intersection.And(intersection, v2.i)
	}

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
	v2.g = v.g

	return
}

func (v Value) BigInt() *big.Int {
	return v.i
}

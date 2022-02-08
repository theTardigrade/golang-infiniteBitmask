package infiniteBitmask

type Pair struct {
	name  string
	value *Value
}

func (p *Pair) Name() string {
	return p.name
}

func (p *Pair) Value() *Value {
	return p.value
}

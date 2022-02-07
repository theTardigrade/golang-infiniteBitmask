package infiniteBitmask

import (
	"math/big"
	"sync"
)

type Generator struct {
	valueCurrent Value
	valuesByName map[string]Value
	mutex        sync.RWMutex
}

func New() (g *Generator) {
	g = &Generator{
		valuesByName: make(map[string]Value),
	}

	g.valueCurrent.i = big.NewInt(0)

	return
}

func (g *Generator) Names() (names []string) {
	defer g.mutex.RUnlock()
	g.mutex.RLock()

	names = make([]string, len(g.valuesByName))

	var i int
	for n := range g.valuesByName {
		names[i] = n
		i++
	}

	return
}

func (g *Generator) Values() (values []Value) {
	defer g.mutex.RUnlock()
	g.mutex.RLock()

	values = make([]Value, len(g.valuesByName))

	var i int
	for _, v := range g.valuesByName {
		values[i] = v
		i++
	}

	return
}

func (g *Generator) ValueFromName(name string) (value Value) {
	value, found := g.valueFromNameReadOnly(name)
	if found {
		return
	}

	value = g.valueFromNameReadWrite(name)

	return
}

func (g *Generator) ValueFromNames(names []string) (value Value) {
	nameValues := make([]Value, len(names))

	for i, n := range names {
		nameValues[i] = g.ValueFromName(n)
	}

	i := new(big.Int)

	for _, v := range nameValues {
		i.Or(i, v.i)
	}

	value.i = i

	return
}

func (g *Generator) valueFromNameReadOnly(name string) (value Value, found bool) {
	defer g.mutex.RUnlock()
	g.mutex.RLock()

	if g.valuesByName == nil {
		return
	}

	value, found = g.valuesByName[name]
	if found {
		value = value.Clone()
	} else {
		return
	}

	return
}

func (g *Generator) valueFromNameReadWrite(name string) (value Value) {
	defer g.mutex.Unlock()
	g.mutex.Lock()

	var found bool

	if g.valuesByName != nil {
		value, found = g.valuesByName[name]
	} else {
		g.valuesByName = make(map[string]Value)
	}

	if found {
		value = value.Clone()
	} else {
		if g.valueCurrent.i == nil {
			g.valueCurrent.i = big.NewInt(0)
		}

		value = g.valueCurrent.Clone()

		g.valueCurrent.i.Lsh(g.valueCurrent.i, 1)

		g.valuesByName[name] = value.Clone()
	}

	return
}

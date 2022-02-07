package infiniteBitmask

import (
	"sync"
)

type Generator struct {
	valueCurrent Value
	valuesByName map[string]Value
	mutex        sync.RWMutex
}

const (
	generatorValueInitial = 1
)

func NewGenerator() (g *Generator) {
	g = &Generator{
		valuesByName: make(map[string]Value),
	}

	g.valueCurrent = g.newValue(generatorValueInitial)

	return
}

func (g *Generator) newValue(n int64) (v Value) {
	v = newValue(n, g)

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

	value = g.newValue(0)

	for _, v := range nameValues {
		value.i.Or(value.i, v.i)
	}

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
			g.valueCurrent = g.newValue(generatorValueInitial)
		}

		value = g.valueCurrent.Clone()

		g.valueCurrent.i.Lsh(g.valueCurrent.i, 1)

		g.valuesByName[name] = value.Clone()
	}

	return
}

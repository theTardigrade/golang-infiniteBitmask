package infiniteBitmask

import (
	"sync"
)

type Generator struct {
	valueCurrent *Value
	valuesByName map[string]*Value
	mutex        sync.RWMutex
}

const (
	generatorValueNumberInitial = 1
)

func NewGenerator() (g *Generator) {
	g = &Generator{
		valuesByName: make(map[string]*Value),
	}

	g.valueCurrent = g.newValue(generatorValueNumberInitial)

	return
}

func (g *Generator) newValue(number int64) (v *Value) {
	v = newValue(number, g)

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

func (g *Generator) Values() (values []*Value) {
	defer g.mutex.RUnlock()
	g.mutex.RLock()

	values = make([]*Value, len(g.valuesByName))

	var i int
	for _, v := range g.valuesByName {
		values[i] = v
		i++
	}

	return
}

func (g *Generator) ValueFromName(name string) (value *Value) {
	value, found := g.valueFromNameReadOnly(name)
	if found {
		return
	}

	value = g.valueFromNameReadWrite(name)

	return
}

func (g *Generator) ValueFromNames(names ...string) (value *Value) {
	value = g.newValue(0)

	for _, n := range names {
		v := g.ValueFromName(n)

		value.number.Or(value.number, v.number)
	}

	return
}

func (g *Generator) valueFromNameReadOnly(name string) (value *Value, found bool) {
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

func (g *Generator) valueFromNameReadWrite(name string) (value *Value) {
	defer g.mutex.Unlock()
	g.mutex.Lock()

	var found bool

	if g.valuesByName != nil {
		value, found = g.valuesByName[name]
	} else {
		g.valuesByName = make(map[string]*Value)
	}

	if found {
		value = value.Clone()
	} else {
		if g.valueCurrent == nil || g.valueCurrent.number == nil {
			g.valueCurrent = g.newValue(generatorValueNumberInitial)
		}

		value = g.valueCurrent.Clone()

		g.valueCurrent.number.Lsh(g.valueCurrent.number, 1)

		g.valuesByName[name] = value.Clone()
	}

	return
}

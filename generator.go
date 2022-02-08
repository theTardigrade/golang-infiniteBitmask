package infiniteBitmask

import (
	"sort"
	"strings"
	"sync"
)

type Generator struct {
	inner       generatorInner
	innerInited bool
}

type generatorInner struct {
	valueCurrent *Value
	valuesByName map[string]*Value
	mutex        sync.RWMutex
}

const (
	generatorValueNumberInitial = 1
)

func NewGenerator() (g *Generator) {
	g = &Generator{}

	g.initInner()

	return
}

func (g *Generator) initInner() {
	g.inner.valuesByName = make(map[string]*Value)
	g.inner.valueCurrent = g.newValue(generatorValueNumberInitial)
	g.innerInited = true
}

func (g *Generator) newValue(number uint8) (v *Value) {
	v = newValue(number, g)

	return
}

func (g *Generator) read(handler func()) {
	if g == nil {
		return
	}

	if !g.innerInited {
		g.initInner()
	}

	defer g.inner.mutex.RUnlock()
	g.inner.mutex.RLock()

	handler()
}

func (g *Generator) write(handler func()) {
	if g == nil {
		return
	}

	if !g.innerInited {
		g.initInner()
	}

	defer g.inner.mutex.Unlock()
	g.inner.mutex.Lock()

	handler()
}

func (g *Generator) Names() (names []string) {
	g.read(func() {
		names = make([]string, len(g.inner.valuesByName))

		var i int
		for n := range g.inner.valuesByName {
			names[i] = n
			i++
		}
	})

	return
}

func (g *Generator) Values() (values []*Value) {
	g.read(func() {
		values = make([]*Value, len(g.inner.valuesByName))

		var i int
		for _, v := range g.inner.valuesByName {
			values[i] = v
			i++
		}
	})

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

		value.inner.number.Or(value.inner.number, v.inner.number)
	}

	return
}

func (g *Generator) valueFromNameReadOnly(name string) (value *Value, found bool) {
	g.read(func() {
		value, found = g.inner.valuesByName[name]
		if found {
			value = value.Clone()
		}
	})

	return
}

func (g *Generator) valueFromNameReadWrite(name string) (value *Value) {
	g.write(func() {
		var found bool

		value, found = g.inner.valuesByName[name]

		if found {
			value = value.Clone()
		} else {
			valueCurrent := g.inner.valueCurrent
			value = valueCurrent.Clone()

			valueCurrent.inner.number.Lsh(valueCurrent.inner.number, 1)

			g.inner.valuesByName[name] = value.Clone()
		}
	})

	return
}

func (g *Generator) String() (result string) {
	names := g.Names()

	sort.Slice(names, func(i, j int) bool {
		vI := g.ValueFromName(names[i])
		vJ := g.ValueFromName(names[j])

		return vI.Number().Cmp(vJ.Number()) == -1
	})

	var builder strings.Builder

	builder.WriteByte('[')

	var i int
	for _, n := range names {
		if i > 0 {
			builder.WriteByte(',')
		}

		builder.WriteByte('"')
		builder.WriteString(strings.ReplaceAll(n, "\"", "\\\""))
		builder.WriteByte('"')

		i++
	}

	builder.WriteByte(']')

	result = builder.String()

	return
}

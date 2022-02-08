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

func NewGeneratorFromString(input string) (g *Generator) {
	g = NewGenerator()

	g.loadString(input)

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

func (g *Generator) Pairs() (pairs []*Pair) {
	g.read(func() {
		pairs = make([]*Pair, len(g.inner.valuesByName))

		var i int
		for n, v := range g.inner.valuesByName {
			pairs[i] = newPair(n, v)
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

func (g *Generator) loadString(input string) (err error) {
	inputLen := len(input)

	if inputLen < 2 {
		err = ErrLoadString
		return
	}

	if input[0] != '[' || input[inputLen-1] != ']' {
		err = ErrLoadString
		return
	}

	inputSplit := make([]string, 0, inputLen/3+1)

	{
		var currInputSplitBuilder strings.Builder
		var prevRune rune

		for _, r := range input[1 : inputLen-1] {
			if r == ',' && prevRune != '\\' {
				if currInputSplitBuilder.Len() < 2 {
					err = ErrLoadString
					return
				}

				inputSplit = append(inputSplit, currInputSplitBuilder.String())
				currInputSplitBuilder.Reset()
			} else {
				currInputSplitBuilder.WriteRune(r)
			}

			prevRune = r
		}

		if l := currInputSplitBuilder.Len(); l > 0 {
			if l < 2 {
				err = ErrLoadString
				return
			}

			inputSplit = append(inputSplit, currInputSplitBuilder.String())
		}
	}

	valueCurrent := g.inner.valueCurrent

	for _, n := range inputSplit {
		nLen := len(n)

		if n[0] != '"' || n[nLen-1] != '"' {
			err = ErrLoadString
			return
		}

		var nameBuilder strings.Builder
		var prevRune rune

		for _, r := range n[1 : nLen-1] {
			switch r {
			case '\\', '"', ',':
				if prevRune == '\\' {
					nameBuilder.WriteByte(byte(r))
				}
			default:
				nameBuilder.WriteRune(r)
			}

			prevRune = r
		}

		name := nameBuilder.String()

		g.inner.valuesByName[name] = valueCurrent.Clone()

		valueCurrent.inner.number.Lsh(valueCurrent.inner.number, 1)
	}

	return
}

func (g *Generator) String() (result string) {
	pairs := g.Pairs()

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].inner.value.inner.number.Cmp(pairs[j].inner.value.inner.number) == -1
	})

	var builder strings.Builder

	builder.WriteByte('[')

	var i int
	for _, p := range pairs {
		if i > 0 {
			builder.WriteByte(',')
		}

		var nameBuilder strings.Builder

		for _, r := range p.inner.name {
			switch r {
			case '\\', '"', ',':
				nameBuilder.WriteByte('\\')
				nameBuilder.WriteByte(byte(r))
			default:
				nameBuilder.WriteRune(r)
			}
		}

		name := nameBuilder.String()

		builder.WriteByte('"')
		builder.WriteString(name)
		builder.WriteByte('"')

		i++
	}

	builder.WriteByte(']')

	result = builder.String()

	return
}

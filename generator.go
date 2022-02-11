package infiniteBitmask

import (
	"math/big"
	"sort"
	"strings"
	"sync"
)

type Generator struct {
	inner       generatorInner
	innerInited bool
	mutex       sync.RWMutex
}

type generatorInner struct {
	valueCurrent *Value
	valuesByName map[string]*Value
}

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
	g.inner.valueCurrent = g.newValue(bigOne)
	g.innerInited = true
}

func (g *Generator) newValue(number *big.Int) (v *Value) {
	v = newValue(number, g)

	return
}

func (g *Generator) read(handler func()) {
	if g == nil {
		panic(ErrPointerNil)
	}

	var shouldWrite bool

	func() {
		defer g.mutex.RUnlock()
		g.mutex.RLock()

		if !g.innerInited {
			shouldWrite = true
			return
		}

		handler()
	}()

	if shouldWrite {
		g.write(handler)
	}
}

func (g *Generator) write(handler func()) {
	if g == nil {
		panic(ErrPointerNil)
	}

	defer g.mutex.Unlock()
	g.mutex.Lock()

	if !g.innerInited {
		g.initInner()
	}

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
			values[i] = v.Clone()
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
			pairs[i] = newPair(n, v.Clone())
			i++
		}
	})

	return
}

func (g *Generator) Clone() (g2 *Generator) {
	g2 = NewGenerator()

	g.read(func() {
		for n, v := range g.inner.valuesByName {
			g2.inner.valuesByName[n] = g2.newValue(v.inner.number)
		}

		g2.inner.valueCurrent = g2.newValue(g.inner.valueCurrent.inner.number)
	})

	return
}

func (g *Generator) Len() (result int) {
	g.read(func() {
		result = len(g.inner.valuesByName)
	})

	return
}

func (g *Generator) Equal(g2 *Generator) (result bool) {
	g.read(func() {
		g2.read(func() {
			if len(g.inner.valuesByName) != len(g2.inner.valuesByName) {
				return
			}

			if !g.inner.valueCurrent.equal(g2.inner.valueCurrent, false) {
				return
			}

			for gName, gValue := range g.inner.valuesByName {
				g2Value, g2NameFound := g2.inner.valuesByName[gName]
				if !g2NameFound {
					return
				}

				if !gValue.equal(g2Value, false) {
					return
				}
			}

			result = true
		})
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
	g.read(func() {
		value = g.newValue(nil)
	})

	valueNumber := value.inner.number

	for _, n := range names {
		v := g.ValueFromName(n)

		valueNumber.Or(valueNumber, v.inner.number)
	}

	return
}

func (g *Generator) ValueFromAllNames() (value *Value) {
	var valueCurrentNumber *big.Int

	g.read(func() {
		value = g.newValue(nil)

		valueCurrentNumber = g.inner.valueCurrent.Number()
	})

	valueNumber := value.inner.number

	for {
		valueCurrentNumber.Rsh(valueCurrentNumber, 1)

		if valueCurrentNumber.Cmp(bigZero) == 0 {
			break
		}

		valueNumber.Or(valueNumber, valueCurrentNumber)
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

	input = input[1 : inputLen-1]

	inputSplitCap := inputLen / 3
	if inputSplitCap < 1 {
		inputSplitCap = 1
	}

	inputSplit := make([]string, 0, inputSplitCap)

	{
		var currInputSplitBuilder strings.Builder
		var prevRune rune

		for _, r := range input {
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

		n = n[1 : nLen-1]

		var nameBuilder strings.Builder
		var prevRune rune

		for _, r := range n {
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

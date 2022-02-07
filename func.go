package infiniteBitmask

import (
	"math/big"
	"sync"
)

var (
	valueCurrent = Value{i: big.NewInt(0)}
	valuesByName = make(map[string]Value)
	mutex        sync.RWMutex
)

func Names() (names []string) {
	defer mutex.RUnlock()
	mutex.RLock()

	names = make([]string, len(valuesByName))

	var i int
	for n := range valuesByName {
		names[i] = n
		i++
	}

	return
}

func Values() (values []Value) {
	defer mutex.RUnlock()
	mutex.RLock()

	values = make([]Value, len(valuesByName))

	var i int
	for _, v := range valuesByName {
		values[i] = v
		i++
	}

	return
}

func ValueFromName(name string) (value Value) {
	value, found := valueFromNameReadOnly(name)
	if found {
		return
	}

	value = valueFromNameReadWrite(name)

	return
}

func ValueFromNames(names []string) (value Value) {
	nameValues := make([]Value, len(names))

	for i, n := range names {
		nameValues[i] = ValueFromName(n)
	}

	i := new(big.Int)

	for _, v := range nameValues {
		i.Or(i, v.i)
	}

	value.i = i

	return
}

func valueFromNameReadOnly(name string) (value Value, found bool) {
	defer mutex.RUnlock()
	mutex.RLock()

	value, found = valuesByName[name]
	if found {
		value = value.Clone()
	} else {
		return
	}

	return
}

func valueFromNameReadWrite(name string) (value Value) {
	defer mutex.Unlock()
	mutex.Lock()

	var found bool

	value, found = valuesByName[name]
	if found {
		value = value.Clone()
	} else {
		value = valueCurrent.Clone()
		found = true

		valueCurrent.i.Lsh(valueCurrent.i, 1)

		valuesByName[name] = value.Clone()
	}

	return
}

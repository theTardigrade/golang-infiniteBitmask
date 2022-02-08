package infiniteBitmask

type Entry struct {
	name      string
	value     *Value
	generator *Generator
}

func (e *Entry) Name() string {
	return e.name
}

func (e *Entry) Value() *Value {
	return e.value
}

func (e *Entry) Generator() *Generator {
	return e.generator
}

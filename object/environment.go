package object

// Environment
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{Store: s}
}

type Environment struct {
	Store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.Store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.Store[name] = val
	return val
}

func (e *Environment) Update(name string, val Object) (Object, bool) {
	_, ok := e.Store[name]
	if ok {
		e.Store[name] = val
		return val, true
	}
	if e.outer != nil {
		return e.outer.Update(name, val)
	}
	return nil, false
}

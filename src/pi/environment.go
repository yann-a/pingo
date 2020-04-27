package pi

import (
	"sync/atomic"
	"unsafe"
)

/**** Concrete values ****/
type Value interface {
	isValue()
}

// channels are values (they can be sent through channels)
type Channel chan Value

func (c Channel) isValue() {}

// type constant is defined in eval.go and is also a terminal
func (c Constant) isValue() {}

type Vpair struct {
	V1 Value
	V2 Value
}

func (c Vpair) isValue() {}

// nothing is defined in eval.go and is also a terminal
func (c Nothing) isValue() {}

/**** Environment ****/
type env struct {
	name  Variable
	value Value
	next  *env
}

/**** Methods ****/
// We append the value at the front of the environment
func (e *env) set_value(x Variable, v Value) *env {
	return &env{x, v, e}
}

func (e *env) get_value(x Variable) Value {
	if e.name == x {
		return e.value
	}

	// If we reach the end of the environment
	if e.next == nil {
		channel := make(Channel)

		// On essaye de mettre à jour le pointeur de fin
		unsafePointer := (*unsafe.Pointer)(unsafe.Pointer(&e.next))
		if atomic.CompareAndSwapPointer(unsafePointer, nil, unsafe.Pointer(&env{x, channel, nil})) {
			return channel
		}

		close(channel) // in case of failure
	}

	// Otherwise we dive deeper in the environment
	return e.next.get_value(x)
}

func (e *env) set_from_pattern(p Terminal, val Value) *env {
	switch pattern := p.(type) {
	case Variable:
		return e.set_value(pattern, val)
	case Pair:
		pair := val.(Vpair)
		env := e.set_value(pattern.V1.(Variable), pair.V1)
		return env.set_value(pattern.V2.(Variable), pair.V2)
	case Nothing:
		return e
	default:
		panic("Not a value provided")
	}
}

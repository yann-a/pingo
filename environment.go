package main

import (
  "sync/atomic"
  "unsafe"
)

/**** Concrete values ****/
type value interface {
	isValue()
}

// channels are values (they can be sent through channels)
type channel chan value

func (c channel) isValue() {}

// type constant is defined in eval.go and is also a terminal
func (c constant) isValue() {}

type vpair struct {
	v1 value
	v2 value
}

func (c vpair) isValue() {}

/**** Environment ****/
type env struct {
	name  variable
	value value
	next  *env
}

/**** Methods ****/
// We append the value at the front of the environment
func (e *env) set_value(x variable, v value) *env {
	return &env{x, v, e}
}

func (e *env) get_value(x variable) value {
	if e.name == x {
		return e.value
	}

	// If we reach the end of the environment
	if e.next == nil {
		channel := make(channel)

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

func (e *env) set_from_pattern(p terminal, val value) *env {
	switch pattern := p.(type) {
	case variable:
		return e.set_value(pattern, val)
	case pair:
		pair := val.(vpair)
		env := e.set_value(pattern.v1.(variable), pair.v1)
		return env.set_value(pattern.v2.(variable), pair.v2)
	default:
		panic("Not a value provided")
	}
}

package typing

import (
	"sync/atomic"
	"unsafe"
  "pingo/src/pi"
)

func createRepr(t Type) *Chain {
	var f Chain = Repr{t}

	return &f
}

/**** Environment ****/
type env struct {
	name  pi.Variable
	value *Chain
	next  *env
}

/**** Methods ****/
// We append the value at the front of the environment
func (e *env) set_value(x pi.Variable, v *Chain) *env {
	return &env{x, v, e}
}

func (e *env) get_value(x pi.Variable) *Chain {
	if e.name == x {
		return e.value
	}

	// If we reach the end of the environment
	if e.next == nil {
    chain := createRepr(Variable{}) // new type variable

		// On essaye de mettre à jour le pointeur de fin
		unsafePointer := (*unsafe.Pointer)(unsafe.Pointer(&e.next))
		if atomic.CompareAndSwapPointer(unsafePointer, nil, unsafe.Pointer(&env{x, chain, nil})) {
			return chain
		}
	}

	// Otherwise we dive deeper in the environment
	return e.next.get_value(x)
}

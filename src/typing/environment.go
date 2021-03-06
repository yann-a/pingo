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
func (e *env) privatize(x pi.Variable, chanType pi.ChanType) *env {
	return &env{x, createRepr(Channel{createRepr(Variable{}), chanType}), e}
}

func (e *env) createVariable(x pi.Variable) *env {
	return &env{x, createRepr(Variable{}), e}
}

func (e *env) get_type(x pi.Variable) *Chain {
	if e.name == x {
		return e.value
	}

	// If we reach the end of the environment, it means we get a channel
	if e.next == nil {
    chain := createRepr(Channel{createRepr(Variable{}), pi.FunChan}) // new type variable

		// On essaye de mettre à jour le pointeur de fin
		unsafePointer := (*unsafe.Pointer)(unsafe.Pointer(&e.next))
		if atomic.CompareAndSwapPointer(unsafePointer, nil, unsafe.Pointer(&env{x, chain, nil})) {
			return chain
		}
	}

	// Otherwise we dive deeper in the environment
	return e.next.get_type(x)
}

func (e *env) type_from_pattern(p pi.Terminal) (*env, *Chain) {
	switch pattern := p.(type) {
	case pi.Variable:
		env := e.createVariable(pattern)

		return env, env.get_type(pattern)

	case pi.Pair:
		env1 := e.createVariable(pattern.V1.(pi.Variable))
		env2 := env1.createVariable(pattern.V2.(pi.Variable))

		return env2, createRepr(Pair{env2.get_type(pattern.V1.(pi.Variable)), env2.get_type(pattern.V2.(pi.Variable))})

	case pi.Nothing:
		return e, createRepr(Void{})
	default:
		panic("Not a pattern provided")
	}
}

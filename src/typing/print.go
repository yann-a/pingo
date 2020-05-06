package typing

import "fmt"

// Types
func (e Int) String() string {
	return "int"
}

func (e Channel) String() string {
	return fmt.Sprintf("%v chan", (*evalRepr(e.Inner)).(Repr).Type)
}

func (e Pair) String() string {
	return fmt.Sprintf("(%v, %v)", e.V1, e.V2)
}

func (e Void) String() string {
	return "void"
}

func (e Variable) String() string {
	return "variable"
}

package typing

import "fmt"

var nextToAssign = 0
var names = make(map[Variable]string)
func nextName() string {
	nextToAssign += 1

	value := nextToAssign
	type_name := ""
	if value == 0 {
		type_name = "a"
	} else {
		for value > 0 {
			modulo := value % 26

			offset := 0
			if value<26 { offset = 1 }

			letter_code := modulo + 97 - offset
			type_name = string(letter_code) + type_name
			value = (value - modulo) / 26
		}
	}

	return "'" + type_name
}

func (e Repr) String() string {
	return fmt.Sprintf("%v", e.Type)
}

func (e Link) String() string {
	return fmt.Sprintf("%v", *e.R)
}

// Types
func (e Int) String() string {
	return "int"
}

func (e Channel) String() string {
	return fmt.Sprintf("%v chan", (*evalRepr(e.Inner)).(Repr).Type)
}

func (e Pair) String() string {
	return fmt.Sprintf("(%v, %v)", *e.V1, *e.V2)
}

func (e Void) String() string {
	return "void"
}

func (e Variable) String() string {
	name, ok := names[e]

	if ok {
		return name
	} else {
		newName := nextName()
		names[e] = newName
		return newName
	}
}

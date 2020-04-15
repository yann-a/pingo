package main

import "fmt"

func eval(e expr, envir *env){
	switch v := e.(type) {
		case parallel:
			for _, task := range(v) {
				go eval(task, envir)
			}
		case send:
			fmt.Printf("TODO : implement type %T (%v)\n", v, v)
		case receiveThen:
			fmt.Printf("TODO : implement type %T (%v)\n", v, v)
		case privatize:
			fmt.Printf("TODO : implement type %T (%v)\n", v, v)
		case print:
			ret, ok := interpretTerminal(v.v, envir)
			integer, ok2 := ret.(constant)
			if !ok || !ok2 {
				fmt.Println("Error: a pair or a channel was provided to print.")
				return
			}

			fmt.Printf("%d\n", int(integer))

			eval(v.then, envir)
		case skip:
			break // nothing to do here
		default:
			fmt.Printf("unrecognised type %T\n", v)
	}
}

func interpretTerminal(val terminal, envir *env) (value, bool) {
	switch v := val.(type) {
		case constant:
			return v, true
		case variable:
			return envir.get_value(v), true
		case pair:
			v1, ok1 := interpretTerminal(v.v1, envir)
			v2, ok2 := interpretTerminal(v.v2, envir)

			return vpair{v1, v2}, ok1 && ok2
		default:
			return constant(0), false
	}
}

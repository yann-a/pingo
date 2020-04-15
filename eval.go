package main

import "fmt"

func eval(e expr, envir env){
	switch v := e.(type) {
		case parallel:
			fmt.Printf("TODO : implement type %T (%v)\n", v, v)
		case send:
			fmt.Printf("TODO : implement type %T (%v)\n", v, v)
		case receiveThen:
			fmt.Printf("TODO : implement type %T (%v)\n", v, v)
		case privatize:
			fmt.Printf("TODO : implement type %T (%v)\n", v, v)
		case print:
			fmt.Printf("TODO : implement type %T (%v)\n", v, v)
		case skip:
			fmt.Printf("TODO : implement type %T (%v)\n", v, v)
		default:
			fmt.Printf("unrecognised type %T\n", v)
	}
}

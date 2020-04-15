package main

import (
	"fmt"
	"sync"
)

func eval(e expr, envir *env, wg *sync.WaitGroup){
	defer wg.Done()

	switch v := e.(type) {
		case parallel:
			for _, task := range(v) {
				wg.Add(1)
				go eval(task, envir, wg)
			}


		case send:
			channel := envir.get_value(variable(v.channel)).(channel)
			val, ok := interpretTerminal(v.value, envir)
			if !ok {
				fmt.Println("Error while sending: not a value provided\n")
				return
			}

			wg.Done()
			channel <- val
			wg.Add(1)


		case receiveThen:
			wg.Done() // on réduit de 1 le compte des eval en cours d'exécution tant qu'on attend un message
			message := <- envir.get_value(variable(v.channel)).(channel)
			wg.Add(1) // On réaugmente quand on finit l'écoute

			switch pattern := v.pattern.(type) {
			case variable:
				envir = envir.set_value(pattern, message)
			case pair:
				pair := message.(vpair)
				envir = envir.set_value(pattern.v1.(variable), pair.v1)
				envir = envir.set_value(pattern.v2.(variable), pair.v2)
			}

			wg.Add(1)
			eval(v.then, envir, wg)


		case privatize:
			envir2 := envir.set_value(variable(v.channel), make(channel))

			wg.Add(1)
			eval(v.then, envir2, wg)


		case print:
			ret, ok := interpretTerminal(v.v, envir)
			if !ok {
				fmt.Println("Error while printing: not a value provided\n")
				return
			}

			integer := ret.(constant)

			fmt.Printf("%d\n", int(integer))

			wg.Add(1)
			eval(v.then, envir, wg)


		case skip:
			// nothing to do here

		case choose:
			var c1, c2 string
			var p, p1, p2 terminal
			var t, t1, t2 expr
			var message value

			switch ve := v.e.(type) {
				case receiveThen:
					c1 = ve.channel
					p1 = ve.pattern
					t1 = ve.then
				default:
					fmt.Println("Processes to choose from can only start by reading : not the case of %v\n", ve)
					return
			}
			switch vf := v.f.(type) {
				case receiveThen:
					c2 = vf.channel
					p2 = vf.pattern
					t2 = vf.then
				default:
					fmt.Println("Processes to choose from can only start by reading : not the case of %v\n", vf)
					return
			}

			wg.Done() // on réduit de 1 le compte des eval en cours d'exécution tant qu'on attend un message
			select {
				case message = <-envir.get_value(variable(c1)).(channel):
					p = p1
					t = t1
				case message = <-envir.get_value(variable(c2)).(channel):
					p = p2
					t = t2
			}
			wg.Add(1) // On réaugmente quand on finit l'écoute

			switch pattern := p.(type) {
				case variable:
					envir = envir.set_value(pattern, message)
				case pair:
					pair := message.(vpair)
					envir = envir.set_value(pattern.v1.(variable), pair.v1)
					envir = envir.set_value(pattern.v2.(variable), pair.v2)
			}

			wg.Add(1)
			eval(t, envir, wg)
		default:
			fmt.Printf("unrecognised type %T (%v)\n", v, v)
	}
}


// Transform a terminal expression into a value
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

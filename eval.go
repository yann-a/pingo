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
			val, ok := interpretTerminal(v.value, envir, wg)
			if !ok {
				fmt.Println("Error while sending: not a value provided\n")
				return
			}

			wg.Add(1)
			channel := envir.get_value(variable(v.channel), wg).(channel)
			channel.input <- val


		case receiveThen:
			channel := envir.get_value(variable(v.channel), wg).(channel)
			privateChan := make(chan value) // on ouvre un canal privé pour communiquer avec l'intermédiaire

			channel.request <- privateChan
			message := <- privateChan // on attend sa réponse sur le canal privé

			switch pattern := v.pattern.(type) {
			case variable:
				envir = envir.set_value(pattern, message)
			case pair:
				pair := message.(vpair)
				envir = envir.set_value(pattern.v1.(variable), pair.v1)
				envir = envir.set_value(pattern.v2.(variable), pair.v2)
			}

			eval(v.then, envir, wg)
			wg.Add(1)


		case privatize:
			envir2 := envir.set_value(variable(v.channel), createChannel(wg))

			eval(v.then, envir2, wg)
			wg.Add(1)


		case print:
			ret, ok := interpretTerminal(v.v, envir, wg)
			if !ok {
				fmt.Println("Error while printing: not a value provided\n")
				return
			}

			integer := ret.(constant)

			fmt.Printf("%d\n", int(integer))

			eval(v.then, envir, wg)
			wg.Add(1)


		case skip:
			// nothing to do here


		// case choose:
		// 	var c1, c2 string
		// 	var p, p1, p2 terminal
		// 	var t, t1, t2 expr
		// 	var message value
		//
		// 	switch ve := v.e.(type) {
		// 		case receiveThen:
		// 			c1 = ve.channel
		// 			p1 = ve.pattern
		// 			t1 = ve.then
		// 		default:
		// 			fmt.Println("Processes to choose from can only start by reading : not the case of %v\n", ve)
		// 			return
		// 	}
		// 	switch vf := v.f.(type) {
		// 		case receiveThen:
		// 			c2 = vf.channel
		// 			p2 = vf.pattern
		// 			t2 = vf.then
		// 		default:
		// 			fmt.Println("Processes to choose from can only start by reading : not the case of %v\n", vf)
		// 			return
		// 	}
		//
		// 	wg.Done() // on réduit de 1 le compte des eval en cours d'exécution tant qu'on attend un message
		// 	select {
		// 		case message = <-envir.get_value(variable(c1)).(channel):
		// 			p = p1
		// 			t = t1
		// 		case message = <-envir.get_value(variable(c2)).(channel):
		// 			p = p2
		// 			t = t2
		// 	}
		// 	wg.Add(1) // On réaugmente quand on finit l'écoute
		//
		// 	switch pattern := p.(type) {
		// 		case variable:
		// 			envir = envir.set_value(pattern, message)
		// 		case pair:
		// 			pair := message.(vpair)
		// 			envir = envir.set_value(pattern.v1.(variable), pair.v1)
		// 			envir = envir.set_value(pattern.v2.(variable), pair.v2)
		// 	}
		//
		// 	wg.Add(1)
		// 	eval(t, envir, wg)


		case conditional:
			val_l, ok_l := interpretTerminal(v.e, envir, wg)
			val_r, ok_r := interpretTerminal(v.f, envir, wg)
			if !ok_l || !ok_r {
				fmt.Printf("Error : can't compare non-values expressions (%v and %v)\n", v.e, v.f)
				return
			}

			if v.eq == (val_l == val_r) {
				eval(v.then, envir, wg)
				wg.Add(1)
			}


		case repl:
			channel := envir.get_value(variable(v.channel), wg).(channel)
			privateChan := make(chan value) // on ouvre un canal privé pour communiquer avec l'intermédiaire

			channel.request <- privateChan
			message := <- privateChan // on attend sa réponse sur le canal privé

			switch pattern := v.pattern.(type) {
			case variable:
				envir = envir.set_value(pattern, message)
			case pair:
				pair := message.(vpair)
				envir = envir.set_value(pattern.v1.(variable), pair.v1)
				envir = envir.set_value(pattern.v2.(variable), pair.v2)
			}

			wg.Add(1)
			go eval(v.then, envir, wg)
			eval(v, envir, wg)
			wg.Add(1)

		default:
			fmt.Printf("unrecognised type %T (%v)\n", v, v)
	}
}


// Transform a terminal expression into a value
func interpretTerminal(val terminal, envir *env, wg *sync.WaitGroup) (value, bool) {
	switch v := val.(type) {
		case constant:
			return v, true
		case variable:
			return envir.get_value(v, wg), true
		case pair:
			v1, ok1 := interpretTerminal(v.v1, envir, wg)
			v2, ok2 := interpretTerminal(v.v2, envir, wg)

			return vpair{v1, v2}, ok1 && ok2
		default:
			return constant(0), false
	}
}

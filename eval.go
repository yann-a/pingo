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
				fmt.Println("Error while sending: not a value provided")
				return
			}

			channel.mux.Lock() // on a besoin de prendre le lock sur le channel avant de l'appeler pour savoir si on va être bloqué ou non
			blocking := *channel.counter > 0
			*channel.counter -= 1
			if blocking { // Aucun receveur attend, le channel va bloqué une goroutine, on réduit le nb de goroutine qui s'exécute
				wg.Done()
			}
			channel.mux.Unlock()

			channel.ch <- val
			if blocking {
				wg.Add(1)
			}


		case receiveThen:
			channel := envir.get_value(variable(v.channel)).(channel)

			channel.mux.Lock()
			blocking := *channel.counter < 0
			*channel.counter += 1
			if blocking { // Aucun envoyeur attend, le channel va bloqué une goroutine, on réduit le nb de goroutine qui s'exécute
				wg.Done()
			}
			channel.mux.Unlock()

			message := <- channel.ch
			if blocking { // On réaugmente quand on finit l'écoute
				wg.Add(1)
			}

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
			envir2 := envir.set_value(variable(v.channel), createChannel())

			wg.Add(1)
			eval(v.then, envir2, wg)


		case print:
			ret, ok := interpretTerminal(v.v, envir)
			if !ok {
				fmt.Println("Error while printing: not a value provided")
				return
			}

			integer := ret.(constant)

			fmt.Printf("%d\n", int(integer))

			wg.Add(1)
			eval(v.then, envir, wg)


		case skip:
			// nothing to do here


		default:
			fmt.Printf("unrecognised type %T\n", v)
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

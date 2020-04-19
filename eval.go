package main

import (
	"fmt"
	"sync"
)

func eval(e expr, envir *env, wg *sync.WaitGroup) {
	defer wg.Done() // We make sure we mark the process as finished at the end

	switch v := e.(type) { // We do some kind of pattern matching on the expression
	case parallel:
		for _, task := range v {
			wg.Add(1)
			// We run each task in parallel in separate goroutines
			go eval(task, envir, wg)
		}

	case send:
		val, ok := interpretTerminal(v.value, envir)
		if !ok {
			fmt.Println("Error while sending: not a value provided\n")
			return
		}

		channel := envir.get_value(variable(v.channel)).(channel)
		channel <- val

	case receiveThen:
		channel := envir.get_value(variable(v.channel)).(channel)
		message := <-channel

		envir = envir.set_from_pattern(v.pattern, message)

		eval(v.then, envir, wg)
		wg.Add(1)

	case privatize:
		envir2 := envir.set_value(variable(v.channel), make(channel, 100))

		eval(v.then, envir2, wg)
		wg.Add(1)

	case print:
		ret, ok := interpretTerminal(v.v, envir)
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

	case choose:
		channel1 := envir.get_value(variable(v.e.channel)).(channel)
		channel2 := envir.get_value(variable(v.f.channel)).(channel)

		select {
		case val := <-channel1:
			envir = envir.set_from_pattern(v.e.pattern, val)
			eval(v.e.then, envir, wg)
			wg.Add(1)

		case val := <-channel2:
			envir = envir.set_from_pattern(v.f.pattern, val)
			eval(v.f.then, envir, wg)
			wg.Add(1)
		}

	case conditional:
		val_l, ok_l := interpretTerminal(v.e, envir)
		val_r, ok_r := interpretTerminal(v.f, envir)
		if !ok_l || !ok_r {
			fmt.Printf("Error : can't compare non-values expressions (%v and %v)\n", v.e, v.f)
			return
		}

		if v.eq == (val_l == val_r) {
			eval(v.then, envir, wg)
			wg.Add(1)
		}

	case repl:
		channel := envir.get_value(variable(v.channel)).(channel)
		message := <-channel

		subenvir := envir.set_from_pattern(v.pattern, message)

		wg.Add(1)
		go eval(v.then, subenvir, wg)
		eval(v, envir, wg)
		wg.Add(1)

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

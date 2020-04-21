package pi

import (
	"fmt"
	"sync"
)

func eval(e Expr, envir *env, wg *sync.WaitGroup) {
	defer wg.Done() // We make sure we mark the process as finished at the end

	switch v := e.(type) { // We do some kind of pattern matching on the expression
	case Parallel:
		for _, task := range v {
			wg.Add(1)
			// We run each task in parallel in separate goroutines
			go eval(task, envir, wg)
		}

	case Send:
		val, ok := interpretTerminal(v.Value, envir)
		if !ok {
			fmt.Println("Error while sending: not a value provided\n")
			return
		}

		wg.Add(1) // We increment the counter to take into account the goroutine that will receive this message
		channel := envir.get_value(Variable(v.Channel)).(Channel)
		channel <- val

	case ReceiveThen:
		channel := envir.get_value(Variable(v.Channel)).(Channel)

		wg.Done()            // in case the goroutine is paused
		message := <-channel // the sender reincrements wg when sending a message

		envir = envir.set_from_pattern(v.Pattern, message)

		eval(v.Then, envir, wg)
		wg.Add(1)

	case Privatize:
		envir2 := envir.set_value(Variable(v.Channel), make(Channel))

		eval(v.Then, envir2, wg)
		wg.Add(1)

	case Print:
		ret, ok := interpretTerminal(v.V, envir)
		if !ok {
			fmt.Println("Error while printing: not a value provided\n")
			return
		}

		integer := ret.(Constant)

		fmt.Printf("%d\n", int(integer))

		eval(v.Then, envir, wg)
		wg.Add(1)

	case Skip:
		// nothing to do here

	case Choose:
		channel1 := envir.get_value(Variable(v.E.Channel)).(Channel)
		channel2 := envir.get_value(Variable(v.F.Channel)).(Channel)

		wg.Done()
		select {
		case val := <-channel1:
			envir = envir.set_from_pattern(v.E.Pattern, val)
			eval(v.E.Then, envir, wg)
			wg.Add(1)

		case val := <-channel2:
			envir = envir.set_from_pattern(v.F.Pattern, val)
			eval(v.F.Then, envir, wg)
			wg.Add(1)
		}

	case Conditional:
		val_l, ok_l := interpretTerminal(v.E, envir)
		val_r, ok_r := interpretTerminal(v.F, envir)
		if !ok_l || !ok_r {
			fmt.Printf("Error : can't compare non-values expressions (%v and %v)\n", v.E, v.F)
			return
		}

		if v.Eq == (val_l == val_r) {
			eval(v.Then, envir, wg)
			wg.Add(1)
		}

	case Repl:
		channel := envir.get_value(Variable(v.Channel)).(Channel)
		wg.Done()
		message := <-channel

		subenvir := envir.set_from_pattern(v.Pattern, message)

		wg.Add(1)
		go eval(v.Then, subenvir, wg)
		eval(v, envir, wg)
		wg.Add(1)

	default:
		fmt.Printf("unrecognised type %T (%v)\n", v, v)
	}
}

// Transform a terminal expression into a value
func interpretTerminal(val Terminal, envir *env) (Value, bool) {
	switch v := val.(type) {
	case Constant:
		return v, true
	case Variable:
		return envir.get_value(v), true
	case Pair:
		v1, ok1 := interpretTerminal(v.V1, envir)
		v2, ok2 := interpretTerminal(v.V2, envir)

		return Vpair{v1, v2}, ok1 && ok2
	default:
		return Constant(0), false
	}
}

package main

import "sync"

/**** Concrete values ****/
type value interface {
	isValue()
}

// channel defined in channel.go
func (c channel) isValue() {}

// type constant is defined in eval.go and is also a terminal
func (c constant) isValue() {}

type vpair struct {
	v1 value
	v2 value
}

func (c vpair) isValue() {}

/**** Environment ****/
type env struct {
	name  variable
	value value
	next  *env
}

/**** Methods ****/
// We append the value at the front of the environment
func (e *env) set_value(x variable, v value) *env {
	return &env{x, v, e}
}

// The end of the environment cannot be updated by two goroutines at the same time
// We prevent that from happening by using a mutex
var accessToEnd sync.Mutex

func (e *env) get_value(x variable, wg *sync.WaitGroup) value {
	if e == nil {
		channel := createChannel(wg)
		// We add the new channel in the global space by appending it at the end of the environment
		e = &env{x, channel, nil}

		return channel
	}

	if e.name == x {
		return e.value
	}

	// If we reach the end of the environment
	if e.next == nil {
		accessToEnd.Lock()         // To avoid interferences between processes, we lock the mutex
		defer accessToEnd.Unlock() // And make sure it's unlocked once we're done

		if e.next == nil { // If no concurrent access before getting the lock
			channel := createChannel(wg)
			// We add the new channel in the global space by appending it at the end of the environment
			e.next = &env{x, channel, nil}
			return channel
		}
	}

	// Otherwise we dive deeper in the environment
	return e.next.get_value(x, wg)
}

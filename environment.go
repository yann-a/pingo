package main

import "sync"

// Concrete values
type value interface {
  isValue()
}

type channel struct {
  ch chan value
  mux *sync.Mutex
  counter *int // On doit compter le nombre d'envoyeurs / récepteurs en attente sur le canal pour savoir combien on a de goroutine active
               // et pouvoir déterminer quand on peut arrêter le process (tout le monde en attente)
}
func (c channel) isValue() { }
func createChannel() channel {
  counter := 0
  return channel{make(chan value), &sync.Mutex{}, &counter}
}

// type constant is defined in eval.go and is also a terminal
func (c constant) isValue() { }

type vpair struct {
  v1 value
  v2 value
}
func (c vpair) isValue() { }


// Environment management
type env struct {
  name variable
  value value
  next *env
}

func (e *env) set_value(x variable, v value) *env {
  return &env{x, v, e}
}


var accessToEnd sync.Mutex // The end of the environment cannot be updated by two goroutines at the same time
                           // We prevent that from happening by using a mutex

func (e *env) get_value(x variable) value {
  if (e == nil) {
    channel := createChannel()
    e = &env{x, channel, nil} // On ajoute le nouveau channel dans l'espace global en le mettant à la fin de l'environnement

    return channel
  }

  if e.name == x {
    return e.value
  }

  if e.next == nil { // on a atteint la fin de l'environnement
    accessToEnd.Lock() // To avoid interferences between processes, we lock the mutex
    defer accessToEnd.Unlock() // And make sure it's unlocked once we're done

    if e.next == nil { // If no concurrent access before getting the lock
      channel := createChannel()
      e.next = &env{x, channel, nil} // On ajoute le nouveau channel dans l'espace global en le mettant à la fin de l'environnement

      return channel
    }
  }


  // Otherwise we dive deeper in the environment
  return e.next.get_value(x)
}

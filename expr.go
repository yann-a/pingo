package main

// Expressions
type expr interface{}



type parallel []expr

type send struct { // send a value on a given channel
  channel string
  value expr
}

type receiveThen struct { // receive a value then execute a process
  channel string
  received terminal
  then expr
}

type privatize struct { // used to define private channels
  channel string
  then expr
}

type sequence struct {
  first expr
  then expr
}

type print struct {
  v terminal
  then expr
}




// Values
type terminal interface {
  isTerminal()
}


type constant int
func (c constant) isTerminal() {}

type variable string
func (c variable) isTerminal() {}

type pair struct {
  v1 terminal
  v2 terminal
}
func (c pair) isTerminal() {}

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
  received value
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

type print expr




// Values
type value interface {
  isValue()
}


type constant int
func (c constant) isValue() {}

type variable string
func (c variable) isValue() {}

type pair struct {
  v1 value
  v2 value
}
func (c pair) isValue() {}

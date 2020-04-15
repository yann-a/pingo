package main

import "fmt"

// Expressions
type expr interface{
	isExpr()
  String() string
}


type parallel []expr
func (e parallel) isExpr() { }
func (e parallel) String() (st string) {
  for i,v := range e {
    if i>0 { st = st+" | "}
    st = st + v.String()
  }
  return
}

type send struct { // send a value on a given channel
  channel string
  value terminal
}
func (e send) isExpr() { }
func (e send) String() string {
  return fmt.Sprintf("^%v %v", e.channel, e.value)
}

type receiveThen struct { // receive a value then execute a process
  channel string
  pattern terminal
  then expr
}
func (e receiveThen) isExpr() { }
func (e receiveThen) String() string {
  return fmt.Sprintf("%v(%v).%v", e.channel, e.pattern, e.then)
}

type privatize struct { // used to define private channels
  channel string
  then expr
}
func (e privatize) isExpr() { }
func (e privatize) String() string {
  return fmt.Sprintf("(%v).%v", e.channel, e.then)
}

type print struct {
  v terminal
  then expr
}
func (e print) isExpr() { }
func (e print) String() string {
  return fmt.Sprintf("print %v; %v", e.v, e.then)
}

type skip int
func (e skip) isExpr() { }
func (e skip) String() (r string) { return }



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

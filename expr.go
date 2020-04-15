package main

// Expressions
type expr interface{
  isExpr()
}



type parallel []expr
func (c parallel) isExpr() { }

type send struct { // send a value on a given channel
  channel string
  value terminal
}
func (c send) isExpr() { }

type receiveThen struct { // receive a value then execute a process
  channel string
  received terminal
  then expr
}
func (c receiveThen) isExpr() { }

type privatize struct { // used to define private channels
  channel string
  then expr
}
func (c privatize) isExpr() { }

type print struct {
  v terminal
  then expr
}
func (c print) isExpr() { }

type skip int
func (c skip) isExpr() { }




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

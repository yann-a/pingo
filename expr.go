package main


// Expressions
type expr interface{
	isExpr()
  String() string
}


type parallel []expr
func (e parallel) isExpr() { }


type send struct { // send a value on a given channel
  channel string
  value terminal
}
func (e send) isExpr() { }


type receiveThen struct { // receive a value then execute a process
  channel string
  pattern terminal
  then expr
}
func (e receiveThen) isExpr() { }


type privatize struct { // used to define private channels
  channel string
  then expr
}
func (e privatize) isExpr() { }


type print struct {
  v terminal
  then expr
}
func (e print) isExpr() { }


type skip int
func (e skip) isExpr() { }


type choose struct {
  e expr
  f expr
}
func (e choose) isExpr() { }



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

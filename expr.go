package main

/**** Expressions ****/
type expr interface {
	isExpr()        // To indicate a type is an expression
	String() string // To print the code
}

// An array of expressions piped together
type parallel []expr

func (e parallel) isExpr() {}

// send a value on a given channel
type send struct {
	channel string
	value   terminal
}

func (e send) isExpr() {}

// receive a value then execute a process
type receiveThen struct {
	channel string
	pattern terminal
	then    expr
}

func (e receiveThen) isExpr() {}

// used to define private channels
type privatize struct {
	channel string
	then    expr
}

func (e privatize) isExpr() {}

// the print instruction, which takes another expression to continue (can be skip)
type print struct {
	v    terminal
	then expr
}

func (e print) isExpr() {}

// the null instruction
type skip int

func (e skip) isExpr() {}

// non-deterministic select between two "receive and execute"-form expressions
type choose struct {
	e receiveThen
	f receiveThen
}

func (e choose) isExpr() {}

// a conditional expression
type conditional struct {
	e    terminal
	eq   bool // true tests equality, false inequality
	f    terminal
	then expr
}

func (e conditional) isExpr() {}

// replication of an input
type repl struct {
	channel string
	pattern terminal
	then    expr
}

func (e repl) isExpr() {}

/**** Values ****/
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

package pi

/**** Expressions ****/
type Expr interface {
	isExpr()        // To indicate a type is an expression
	String() string // To print the code
}

// An array of expressions piped together
type Parallel []Expr

func (e Parallel) isExpr() {}

// send a value on a given channel
type Send struct {
	Channel string
	Value   Terminal
}

func (e Send) isExpr() {}

// receive a value then execute a process
type ReceiveThen struct {
	Channel string
	Pattern Terminal
	Then    Expr
}

func (e ReceiveThen) isExpr() {}

// used to define private channels
type Privatize struct {
	Channel string
	Then    Expr
}

func (e Privatize) isExpr() {}

// the print instruction, which takes another expression to continue (can be skip)
type Print struct {
	V    Terminal
	Then Expr
}

func (e Print) isExpr() {}

// the null instruction
type Skip int

func (e Skip) isExpr() {}

// non-deterministic select between two "receive and execute"-form expressions
type Choose struct {
	E ReceiveThen
	F ReceiveThen
}

func (e Choose) isExpr() {}

// a conditional expression
type Conditional struct {
	E    Terminal
	Eq   bool // true tests equality, false inequality
	F    Terminal
	Then Expr
}

func (e Conditional) isExpr() {}

// replication of an input
type Repl struct {
	Channel string
	Pattern Terminal
	Then    Expr
}

func (e Repl) isExpr() {}

/**** Values ****/
type Terminal interface {
	isTerminal()
}

type Constant int
func (c Constant) isTerminal() {}

type Variable string
func (c Variable) isTerminal() {}

type Pair struct {
	V1 Terminal
	V2 Terminal
}
func (c Pair) isTerminal() {}

type Add struct {
	V1 Terminal
	V2 Terminal
}
func (c Add) isTerminal() {}

type Sub struct {
	V1 Terminal
	V2 Terminal
}
func (c Sub) isTerminal() {}

type Mul struct {
	V1 Terminal
	V2 Terminal
}
func (c Mul) isTerminal() {}

type Div struct {
	V1 Terminal
	V2 Terminal
}
func (c Div) isTerminal() {}

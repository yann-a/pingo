package lambda

type Lambda interface {
	isLambda() // To indicate a type is an expression
	//String() string // To print the code
}

type Lconst int

func (l Lconst) isLambda() {}

type Lvar string

func (l Lvar) isLambda() {}

type Lfun struct {
	Arg Lvar
	Exp Lambda
}

func (l Lfun) isLambda() {}

type Lapp struct {
	Fun Lambda
	Exp Lambda
}

func (l Lapp) isLambda() {}

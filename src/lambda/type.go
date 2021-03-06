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

type Add struct {
	L Lambda
	R Lambda
}

func (l Add) isLambda() {}

type Sub struct {
	L Lambda
	R Lambda
}

func (l Sub) isLambda() {}

type Mult struct {
	L Lambda
	R Lambda
}

func (l Mult) isLambda() {}

type Div struct {
	L Lambda
	R Lambda
}

func (l Div) isLambda() {}

type Read struct {
	Var  Lvar
	Ref  Lambda
	Then Lambda
}

func (l Read) isLambda() {}

type Write struct {
	Ref  Lambda
	Val  Lambda
	Then Lambda
}

func (l Write) isLambda() {}

type Swap struct {
	Var  Lvar
	Ref  Lambda
	Val  Lambda
	Then Lambda
}

func (l Swap) isLambda() {}

// Permet d'instancier une réf
type New struct {
	Var   Lvar
	Value Lambda
	Then  Lambda
}

func (l New) isLambda() {}

type Deref struct {
	Name Lambda
}

func (l Deref) isLambda() {}

package lambda

var ShutUp = "I'm just here to make go shut up"

type lambda interface {
	isLambda()        // To indicate a type is an expression
	//String() string // To print the code
}

type lconst int
func (l lconst) isLambda() { }

type lvar string
func (l lvar) isLambda() { }

type lfun struct {
	arg lvar
	exp lambda
}
func (l lfun) isLambda() { }

type lapp struct {
	fun lambda
	exp lambda
}
func (l lapp) isLambda() { }
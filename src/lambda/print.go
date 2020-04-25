package lambda

import "fmt"

/*func (e lconst) String() string {
	return fmt.Sprintf("%d", e)
}

func (e lvar) String() string {
	return fmt.Sprintf("%s", e)
}*/

func (e Lfun) String() string {
	return fmt.Sprintf("(fun %v -> %v)", e.Arg, e.Exp)
}

func (e Lapp) String() string {
	return fmt.Sprintf("(%v %v)", e.Fun, e.Exp)
}

func (e Add) String() string {
	return fmt.Sprintf("(%v + %v)", e.L, e.R)
}

func (e Sub) String() string {
	return fmt.Sprintf("(%v - %v)", e.L, e.R)
}

func (e Mult) String() string {
	return fmt.Sprintf("(%v * %v)", e.L, e.R)
}

func (e Div) String() string {
	return fmt.Sprintf("(%v / %v)", e.L, e.R)
}

func (e Read) String() string {
	return fmt.Sprintf("(%v<-%v ; %v)", e.Var, e.Ref, e.Then)
}

func (e Write) String() string {
	return fmt.Sprintf("(%v:=%v ; %v)", e.Ref, e.Val, e.Then)
}

func (e Swap) String() string {
	return fmt.Sprintf("(%v<-%v:=%v ; %v", e.Var, e.Ref, e.Val, e.Then)
}

func (e New) String() string {
	return fmt.Sprintf("new %v = %v in %v", e.Var, e.Value, e.Then)
}

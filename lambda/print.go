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
package lambda

import "fmt"

/*func (e lconst) String() string {
	return fmt.Sprintf("%d", e)
}

func (e lvar) String() string {
	return fmt.Sprintf("%s", e)
}*/

func (e lfun) String() string {
	return fmt.Sprintf("(fun %v -> %v)", e.arg, e.exp)
}

func (e lapp) String() string {
	return fmt.Sprintf("(%v %v)", e.fun, e.exp)
}
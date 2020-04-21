package pi

import "fmt"

/* We implement the String() method from fmt, used by the print functions, for the expressions */
func (e Parallel) String() (st string) {
	st = "("
	for i, v := range e {
		if i > 0 { // We put a pipe before every expressions but the first
			st = st + " | "
		}
		st = st + v.String()
	}
	st += ")"
	return
}

func (e Send) String() string {
	return fmt.Sprintf("(^%v %v)", e.Channel, e.Value)
}

func (e ReceiveThen) String() string {
	return fmt.Sprintf("(%v(%v).%v)", e.Channel, e.Pattern, e.Then)
}

func (e Privatize) String() string {
	return fmt.Sprintf("((%v).%v)", e.Channel, e.Then)
}

func (e Print) String() string {
	return fmt.Sprintf("(print %v; %v)", e.V, e.Then)
}

func (e Skip) String() string {
	return "skip"
}

func (e Choose) String() string {
	return fmt.Sprintf("(%v + %v)", e.E, e.F)
}

func (e Conditional) String() string {
	var condS string
	if e.Eq {
		condS = "="
	} else {
		condS = "!="
	}

	return fmt.Sprintf("([%v %s %v]%v)", e.E, condS, e.F, e.Then)
}

func (e Repl) String() string {
	return fmt.Sprintf("(!%v(%v).%v)", e.Channel, e.Pattern, e.Then)
}

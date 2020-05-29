package translate

import (
	"fmt"
	"pingo/src/lambda"
)

func ocamlToLambda(lexpr lambda.Lambda, n *int) lambda.Lambda {
	switch v := lexpr.(type) {
		case lambda.Lfun:
			return lambda.Lfun{v.Arg, ocamlToLambda(v.Exp, n)}
		case lambda.Lapp:
			return lambda.Lapp{ocamlToLambda(v.Fun, n), ocamlToLambda(v.Exp, n)}
		case lambda.Add:
			return lambda.Add{ocamlToLambda(v.L, n), ocamlToLambda(v.R, n)}
		case lambda.Sub:
			return lambda.Sub{ocamlToLambda(v.L, n), ocamlToLambda(v.R, n)}
		case lambda.Mult:
			return lambda.Mult{ocamlToLambda(v.L, n), ocamlToLambda(v.R, n)}
		case lambda.Div:
			return lambda.Div{ocamlToLambda(v.L, n), ocamlToLambda(v.R, n)}
		case lambda.Read:
			n0 := *n
			*n = n0+1
			return lambda.Lapp{lambda.Lfun{lambda.Lvar(fmt.Sprintf("thisvarshouldnotbeused%d", n0)), lambda.Read{v.Var, lambda.Lvar(fmt.Sprintf("thisvarshouldnotbeused%d", n0)), ocamlToLambda(v.Then, n)}}, ocamlToLambda(v.Ref, n)}
		case lambda.Write:
			n0 := *n
			*n = n0+1
			return lambda.Lapp{lambda.Lfun{lambda.Lvar(fmt.Sprintf("thisvarshouldnotbeused%d", n0)), lambda.Write{lambda.Lvar(fmt.Sprintf("thisvarshouldnotbeused%d", n0)), v.Val, ocamlToLambda(v.Then, n)}}, ocamlToLambda(v.Ref, n)}
		case lambda.Swap:
			n0 := *n
			*n = n0+1
			return lambda.Lapp{lambda.Lfun{lambda.Lvar(fmt.Sprintf("thisvarshouldnotbeused%d", n0)), lambda.Swap{v.Var, lambda.Lvar(fmt.Sprintf("thisvarshouldnotbeused%d", n0)), v.Val, ocamlToLambda(v.Then, n)}}, ocamlToLambda(v.Ref, n)}
		case lambda.New:
			return lambda.New{v.Var, ocamlToLambda(v.Value, n), ocamlToLambda(v.Then, n)}
		case lambda.Deref:
			return lambda.Deref{ocamlToLambda(v.Name, n)}
		default:
			return v
	}
}
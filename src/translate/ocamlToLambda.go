package translate

import "pingo/src/lambda"


// Performs the caml -> lambdaref translation
func ocamlToLambda(lexpr lambda.Lambda) lambda.Lambda {
	switch v := lexpr.(type) {
		case lambda.Lfun:
			return lambda.Lfun{v.Arg, ocamlToLambda(v.Exp)}
		case lambda.Lapp:
			return lambda.Lapp{ocamlToLambda(v.Fun), ocamlToLambda(v.Exp)}
		case lambda.Add:
			return lambda.Add{ocamlToLambda(v.L), ocamlToLambda(v.R)}
		case lambda.Sub:
			return lambda.Sub{ocamlToLambda(v.L), ocamlToLambda(v.R)}
		case lambda.Mult:
			return lambda.Mult{ocamlToLambda(v.L), ocamlToLambda(v.R)}
		case lambda.Div:
			return lambda.Div{ocamlToLambda(v.L), ocamlToLambda(v.R)}
		case lambda.Read:
			return lambda.Lapp{lambda.Lfun{lambda.Lvar("thisvarshouldnotbeused"), lambda.Read{v.Var, lambda.Lvar("thisvarshouldnotbeused"), ocamlToLambda(v.Then)}}, ocamlToLambda(v.Ref)}
		case lambda.Write:
			return lambda.Lapp{lambda.Lfun{lambda.Lvar("thisvarshouldnotbeused"), lambda.Write{lambda.Lvar("thisvarshouldnotbeused"), v.Val, ocamlToLambda(v.Then)}}, ocamlToLambda(v.Ref)}
		case lambda.Swap:
			return lambda.Lapp{lambda.Lfun{lambda.Lvar("thisvarshouldnotbeused"), lambda.Swap{v.Var, lambda.Lvar("thisvarshouldnotbeused"), v.Val, ocamlToLambda(v.Then)}}, ocamlToLambda(v.Ref)}
		case lambda.New:
			return lambda.New{v.Var, ocamlToLambda(v.Value), ocamlToLambda(v.Then)}
		case lambda.Deref:
			return lambda.Deref{ocamlToLambda(v.Name)}
		default:
			return v
	}
}
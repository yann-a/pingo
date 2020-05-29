package translate

import (
	"fmt"
	"pingo/src/lambda"
	"pingo/src/pi"
)

func Translate(lexpr lambda.Lambda, channel int) pi.Expr {
	n := 0
	tr1 := ocamlToLambda(lexpr, &n)
	fmt.Println(tr1)
	translation := innerTranslate(tr1, channel)

	return pi.Parallel{
		translation,

		pi.Repl{ // On définit print comme une fonction lambda usuelle
			"print",
			pi.Pair{pi.Variable("x"), pi.Variable("q")},
			pi.Print{pi.Variable("x"), pi.Send{"q", pi.Variable("x")}},
		},
	}
}

func chanName(channel int) string {
	return fmt.Sprintf("cont%d", channel)
}

func innerTranslate(lexpr lambda.Lambda, channel int) pi.Expr {
	// on détermine des noms frais pour les translate récursifs
	channel1 := (channel + 1) % 3
	channel2 := (channel + 2) % 3

	switch v := lexpr.(type) {
	case lambda.Lconst:
		return pi.Send{chanName(channel), pi.Constant(v)}
	case lambda.Lvar:
		return pi.Send{chanName(channel), pi.Variable(v)}
	case lambda.Lfun:
		// Une fonction lambda est transformée en un canal qui reçoit des paires (argument, canal de retour)
		return pi.Privatize{
			"y",
			pi.FunChan,
			pi.Parallel{
				pi.Send{
					chanName(channel),
					pi.Variable("y"),
				},
				pi.Repl{
					"y",
					pi.Pair{pi.Variable(v.Arg), pi.Variable(chanName(0))},
					innerTranslate(v.Exp, 0),
				},
			},
		}
	case lambda.Lapp:
		return translateArith(
			v.Fun,
			v.Exp,
			func(L, R pi.Terminal) pi.Expr {
				return pi.Send{
					string(L.(pi.Variable)),
					pi.Pair{R, pi.Variable(chanName(channel))},
				}
			},
			channel1,
			channel2,
		)
	case lambda.Add:
		return translateArith(
			v.L,
			v.R,
			func(L, R pi.Terminal) pi.Expr {
				return pi.Send{
					chanName(channel),
					pi.Add{L, R},
				}
			},
			channel1,
			channel2,
		)
	case lambda.Sub:
		return translateArith(
			v.L,
			v.R,
			func(L, R pi.Terminal) pi.Expr {
				return pi.Send{
					chanName(channel),
					pi.Sub{L, R},
				}
			},
			channel1,
			channel2,
		)
	case lambda.Mult:
		return translateArith(
			v.L,
			v.R,
			func(L, R pi.Terminal) pi.Expr {
				return pi.Send{
					chanName(channel),
					pi.Mul{L, R},
				}
			},
			channel1,
			channel2,
		)
	case lambda.Div:
		return translateArith(
			v.L,
			v.R,
			func(L, R pi.Terminal) pi.Expr {
				return pi.Send{
					chanName(channel),
					pi.Div{L, R},
				}
			},
			channel1,
			channel2,
		)
	case lambda.Read:
		return pi.ReceiveThen{
			string(v.Ref.(lambda.Lvar)),
			pi.Variable(v.Var),
			pi.Parallel{
				pi.Send{string(v.Ref.(lambda.Lvar)), pi.Variable(v.Var)},
				innerTranslate(v.Then, channel),
			},
		}
	case lambda.Write:
		return translatePrimitives(pi.Variable("_"), string(v.Ref.(lambda.Lvar)), v.Val, v.Then, channel, channel1)
	case lambda.Swap:
		return translatePrimitives(pi.Variable(v.Var), string(v.Ref.(lambda.Lvar)), v.Val, v.Then, channel, channel1)
	case lambda.New:
		var t pi.Terminal = translateLambdaTerminal(v.Value)
		if t == nil { // Si l'expression écrit n'est pas simple
			t = pi.Variable("value")
		}

		finalExpr := pi.Privatize{ // On réserve la variable de la future réf
			string(v.Var),
			pi.RefChan,
			pi.Parallel{
				pi.Send{string(v.Var), t}, // émission de la valeur sur le canal de la réf
				innerTranslate(v.Then, channel),
			},
		}

		if t == nil {
			finalExpr = pi.Privatize{ // On privatise un canal pour attendre la réception de la valeur écrite
				chanName(channel1),
				pi.FunChan,
				pi.Parallel{
					innerTranslate(v.Value, channel1),
					pi.ReceiveThen{
						chanName(channel1),
						pi.Variable("value"),
						finalExpr,
					},
				},
			}
		}

		return finalExpr
	case lambda.Deref:
		return pi.Privatize{
			chanName(channel1),
			pi.FunChan,
			pi.Parallel{
				innerTranslate(v.Name, channel1),
				pi.ReceiveThen{
					chanName(channel1),
					pi.Variable("refName"),
					pi.ReceiveThen{
						"refName",
						pi.Variable("refContent"),
						pi.Parallel{
							pi.Send{"refName", pi.Variable("refContent")},
							pi.Send{chanName(channel), pi.Variable("refContent")},
						},
					},
				},
			},
		}
	default:
		panic("not supposed to happen")
	}
}

func translateLambdaTerminal(lexpr lambda.Lambda) pi.Terminal {
	switch v := lexpr.(type) {
	case lambda.Lconst:
		return pi.Constant(v)
	case lambda.Lvar:
		return pi.Variable(v)
	default:
		return nil
	}
}

func translateArith(L, R lambda.Lambda, sendExpr func(L, R pi.Terminal) pi.Expr, channel1, channel2 int) pi.Expr {
	t1 := translateLambdaTerminal(L)
	t2 := translateLambdaTerminal(R)

	var rresult pi.Terminal = pi.Variable("rresult")
	if t2 != nil { // R est simple, on peut directement récupérer sa valeur en pi
		rresult = t2
	}

	var lresult pi.Terminal = pi.Variable("lresult")
	if t1 != nil { // R est simple, on peut directement récupérer sa valeur en pi
		lresult = t1
	}

	finalExpr := sendExpr(lresult, rresult)

	if t1 == nil { // L n'est pas simple, on utilise une continuation
		finalExpr = pi.Privatize{
			chanName(channel1),
			pi.FunChan,
			pi.Parallel{
				innerTranslate(L, channel1),
				pi.ReceiveThen{
					chanName(channel1),
					pi.Variable("lresult"),
					finalExpr,
				},
			},
		}
	}

	if t2 == nil { // R n'est pas simple, on utilise une continuation
		finalExpr = pi.Privatize{
			chanName(channel2),
			pi.FunChan,
			pi.Parallel{
				innerTranslate(R, channel2),
				pi.ReceiveThen{
					chanName(channel2),
					pi.Variable("rresult"),
					finalExpr,
				},
			},
		}
	}

	return finalExpr
}

// This is basically a swap, and we us it for swap and write:
// * A write is a swap where we trash the received value
// * A swap is... well, a swap
func translatePrimitives(variable pi.Terminal, ref string, value, then lambda.Lambda, channel, channel1 int) pi.Expr {
	return pi.Privatize{
		chanName(channel1),
		pi.FunChan,
		pi.Parallel{
			innerTranslate(value, channel1),
			pi.ReceiveThen{
				ref,
				variable,
				pi.ReceiveThen{
					chanName(channel1),
					pi.Variable("retrans"),
					pi.Parallel{
						pi.Send{ref, pi.Variable("retrans")},
						innerTranslate(then, channel),
					},
				},
			},
		},
	}
}

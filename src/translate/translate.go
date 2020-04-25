package translate

import (
	"pingo/src/lambda"
	"pingo/src/pi"
)

func Translate(lexpr lambda.Lambda, channel string) pi.Expr {
	freechannel := "q"
	if channel == "q" {
		freechannel = "r"
	}

	translation := innerTranslate(lexpr, freechannel)

	return pi.Parallel{
		translation,

		pi.Repl{ // On définit print comme une fonction lambda usuelle
			"print",
			pi.Pair{pi.Variable("x"), pi.Variable("q")},
			pi.Print{pi.Variable("x"), pi.Send{"q", pi.Variable("x")}},
		},

		// On rajoute un cleaner pour nettoyer les réfs après leur création
		pi.ReceiveThen{
			freechannel,
			pi.Variable("ret"),
			pi.Parallel{
				pi.Repl{
					"refCleaner",
					pi.Variable("a"),
					pi.ReceiveThen{
						"a",
						pi.Variable("x"),
						pi.Skip(0),
					},
				},
				pi.Send{channel, pi.Variable("ret")},
			},
		},
	}
}

func innerTranslate(lexpr lambda.Lambda, channel string) pi.Expr {
	// on détermine des noms frais pour les translate récursifs
	var channel1, channel2 string
	if channel == "q" {
		channel1 = "r"
		channel2 = "s"
	} else if channel == "r" {
		channel1 = "q"
		channel2 = "s"
	} else {
		channel1 = "q"
		channel2 = "r"
	}

	switch v := lexpr.(type) {
	case lambda.Lconst:
		return pi.Send{channel, pi.Constant(v)}
	case lambda.Lvar:
		return pi.Send{channel, pi.Variable(v)}
	case lambda.Lfun:
		// Une fonction lambda est transformée en un canal qui reçoit des paires (argument, canal de retour)
		return pi.Privatize{
			"y",
			pi.Parallel{
				pi.Send{
					channel,
					pi.Variable("y"),
				},
				pi.Repl{
					"y",
					pi.Pair{pi.Variable(v.Arg), pi.Variable("q")},
					innerTranslate(v.Exp, "q"),
				},
			},
		}
	case lambda.Lapp:
		return translateArith(
			v.Fun,
			v.Exp,
			func (L, R pi.Terminal) pi.Expr {
				return pi.Send{
					string(L.(pi.Variable)),
					pi.Pair{R, pi.Variable(channel)},
				}
			},
			channel1,
			channel2,
		)
	case lambda.Add:
		return translateArith(
			v.L,
			v.R,
			func (L, R pi.Terminal) pi.Expr {
				return pi.Send{
					channel,
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
			func (L, R pi.Terminal) pi.Expr {
				return pi.Send{
					channel,
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
			func (L, R pi.Terminal) pi.Expr {
				return pi.Send{
					channel,
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
			func (L, R pi.Terminal) pi.Expr {
				return pi.Send{
					channel,
					pi.Div{L, R},
				}
			},
			channel1,
			channel2,
		)
	case lambda.Read:
		return pi.ReceiveThen{
			string(v.Ref),
			pi.Variable(v.Var),
			pi.Parallel{
				pi.Send{string(v.Ref), pi.Variable(v.Var)},
				innerTranslate(v.Then, channel),
			},
		}
	case lambda.Write:
		return translatePrimitives("trash", string(v.Ref), v.Val, v.Then, channel, channel1)
	case lambda.Swap:
		return translatePrimitives(string(v.Var), string(v.Ref), v.Val, v.Then, channel, channel1)
	case lambda.New:
		var t pi.Terminal = translateLambdaTerminal(v.Value)
		if t == nil { // Si l'expression écrit n'est pas simple
			t = pi.Variable("value")
		}

		finalExpr := pi.Privatize{ // On réserve la variable de la future réf
			string(v.Var),
			pi.Parallel{
				pi.Send{"refCleaner", pi.Variable(v.Var)},
				pi.Send{string(v.Var), t}, // émission de la valeur sur le canal de la réf
				innerTranslate(v.Then, channel),
			},
		}

		if t == nil {
			finalExpr = pi.Privatize{ // On privatise un canal pour attendre la réception de la valeur écrite
				channel1,
				pi.Parallel{
					innerTranslate(v.Value, channel1),
					pi.ReceiveThen{
						channel1,
						pi.Variable("value"),
						finalExpr,
					},
				},
			}
		}

		return finalExpr
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

func translateArith(L, R lambda.Lambda, sendExpr func (L, R pi.Terminal) pi.Expr, channel1, channel2 string) pi.Expr {
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
			channel1,
			pi.Parallel{
				innerTranslate(L, channel1),
				pi.ReceiveThen{
					channel1,
					pi.Variable("lresult"),
					finalExpr,
				},
			},
		}
	}

	if t2 == nil { // R n'est pas simple, on utilise une continuation
		finalExpr = pi.Privatize{
			channel2,
			pi.Parallel{
				innerTranslate(R, channel2),
				pi.ReceiveThen{
					channel2,
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
func translatePrimitives(variable, ref string, value, then lambda.Lambda, channel, channel1 string) pi.Expr {
	return pi.Privatize{
		channel1,
		pi.Parallel{
			innerTranslate(value, channel1),
			pi.ReceiveThen{
				ref,
				pi.Variable(variable),
				pi.Parallel{
					pi.ReceiveThen{
						channel1,
						pi.Variable("retrans"),
						pi.Send{ref, pi.Variable("retrans")},
					},
					innerTranslate(then, channel),
				},
			},
		},
	}
}

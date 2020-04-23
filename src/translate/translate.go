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

		pi.Repl{ // On définit ref
			"ref",
			pi.Pair{pi.Variable("x"), pi.Variable("q")},
			pi.Privatize{
				"a",
				pi.Parallel{
					pi.Send{"a", pi.Variable("x")},
					pi.Send{"refCleaner", pi.Variable("a")},
					pi.Send{"q", pi.Variable("a")},
				},
			},
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
			pi.Send{
				"lresult",
				pi.Pair{pi.Variable("rresult"), pi.Variable(channel)},
			},
			channel1,
			channel2,
		)
	case lambda.Add:
		return translateArith(
			v.L,
			v.R,
			pi.Send{
				channel,
				pi.Add{pi.Variable("lresult"), pi.Variable("rresult")},
			},
			channel1,
			channel2,
		)
	case lambda.Sub:
		return translateArith(
			v.L,
			v.R,
			pi.Send{
				channel,
				pi.Sub{pi.Variable("lresult"), pi.Variable("rresult")},
			},
			channel1,
			channel2,
		)
	case lambda.Mult:
		return translateArith(
			v.L,
			v.R,
			pi.Send{
				channel,
				pi.Mul{pi.Variable("lresult"), pi.Variable("rresult")},
			},
			channel1,
			channel2,
		)
	case lambda.Div:
		return translateArith(
			v.L,
			v.R,
			pi.Send{
				channel,
				pi.Div{pi.Variable("lresult"), pi.Variable("rresult")},
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
	// Since a write is basically a swap where we trash the received value, we could
	// try to factor the two cases but i didn't managed to do it for now
	// (If we do succeed it's probably possible to factor read too)
	case lambda.Write:
		return pi.Privatize{
			channel1,
			pi.ReceiveThen{
				string(v.Ref),
				pi.Variable("trash"),
				pi.Parallel{
					innerTranslate(v.Val, channel1),
					pi.ReceiveThen{
						channel1,
						pi.Variable("retrans"),
						pi.Send{string(v.Ref), pi.Variable("retrans")},
					},
				},
			},
		}
	case lambda.Swap:
		return pi.Privatize{
			channel1,
			pi.ReceiveThen{
				string(v.Ref),
				pi.Variable(v.Var),
				pi.Parallel{
					innerTranslate(v.Val, channel1),
					pi.ReceiveThen{
						channel1,
						pi.Variable("retrans"),
						pi.Send{string(v.Ref), pi.Variable("retrans")},
					},
				},
			},
		}
	default:
		panic("not supposed to happen")
	}
}

func translateArith(L, R lambda.Lambda, sendExpr pi.Expr, channel1, channel2 string) pi.Expr {
	return pi.Privatize{
		channel2,
		pi.Parallel{
			innerTranslate(R, channel2),
			pi.ReceiveThen{
				channel2,
				pi.Variable("rresult"),
				pi.Privatize{
					channel1,
					pi.Parallel{
						innerTranslate(L, channel1),
						pi.ReceiveThen{
							channel1,
							pi.Variable("lresult"),
							sendExpr,
						},
					},
				},
			},
		},
	}
}

package translate

import (
	"pingo/src/lambda"
	"pingo/src/pi"
)

func Translate(lexpr lambda.Lambda, channel string) pi.Expr {
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
					Translate(v.Exp, "q"),
				},
			},
		}
	case lambda.Lapp:
		return translateArith(v.Fun, v.Exp, pi.Send{"lresult", pi.Pair{pi.Variable("rresult"), pi.Variable(channel)}}, channel1, channel2)
	case lambda.Add:
		return translateArith(v.L, v.R, pi.Send{channel, pi.Add{pi.Variable("lresult"), pi.Variable("rresult")}}, channel1, channel2)
	case lambda.Sub:
		return translateArith(v.L, v.R, pi.Send{channel, pi.Sub{pi.Variable("lresult"), pi.Variable("rresult")}}, channel1, channel2)
	case lambda.Mult:
		return translateArith(v.L, v.R, pi.Send{channel, pi.Mul{pi.Variable("lresult"), pi.Variable("rresult")}}, channel1, channel2)
	case lambda.Div:
		return translateArith(v.L, v.R, pi.Send{channel, pi.Div{pi.Variable("lresult"), pi.Variable("rresult")}}, channel1, channel2)
	case lambda.Print:
		return pi.Privatize{channel1, pi.Parallel{Translate(v.L, channel1), pi.ReceiveThen{channel1, pi.Variable("result"), pi.Print{pi.Variable("result"), pi.Send{channel, pi.Variable("result")}}}}}
	default:
		panic("not supposed to happen")
	}
}

func translateArith(L lambda.Lambda, R lambda.Lambda, sendExpr pi.Expr, channel1 string, channel2 string) pi.Expr {
	return pi.Privatize{channel2, pi.Parallel{Translate(R, channel2), pi.ReceiveThen{channel2, pi.Variable("rresult"), pi.Privatize{channel1, pi.Parallel{Translate(R, channel1), pi.ReceiveThen{channel1, pi.Variable("lresult"), sendExpr}}}}}}
}

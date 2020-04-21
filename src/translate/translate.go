package translate

import (
	"pingo/src/lambda"
	"pingo/src/pi"
)

func Translate(lexpr lambda.Lambda, channel string) pi.Expr {
	switch v := lexpr.(type) {
	case lambda.Lconst:
		return pi.Send{channel, pi.Constant(v)}
	case lambda.Lvar:
		return pi.Send{channel, pi.Variable(v)}
	case lambda.Lfun:
		return pi.Privatize{"y", pi.Parallel{pi.Send{channel, pi.Variable("y")}, pi.ReceiveThen{"y", pi.Pair{pi.Variable(v.Arg), pi.Variable("q")}, Translate(v.Exp, "q")}}}
	case lambda.Lapp:
		return pi.Privatize{"q", pi.Parallel{Translate(v.Exp, "q"), pi.ReceiveThen{"q", pi.Variable("v"), pi.Privatize{"r", pi.Parallel{Translate(v.Fun, "r"), pi.ReceiveThen{"r", pi.Variable("f"), pi.Send{"f", pi.Pair{pi.Variable("v"), pi.Variable("p")}}}}}}}}
	default:
		panic("not supposed to happen")
	}
}

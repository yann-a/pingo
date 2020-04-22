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
	case lambda.Add:
		return pi.Privatize{"l", pi.Privatize{"r", pi.Parallel{Translate(v.L, "l"), Translate(v.R, "r"), pi.ReceiveThen{"l", pi.Variable("lresult"), pi.ReceiveThen{"r", pi.Variable("rresult"), pi.Send{"p", pi.Add{pi.Variable("lresult"), pi.Variable("rresult")}}}}}}}
	case lambda.Sub:
		return pi.Privatize{"l", pi.Privatize{"r", pi.Parallel{Translate(v.L, "l"), Translate(v.R, "r"), pi.ReceiveThen{"l", pi.Variable("lresult"), pi.ReceiveThen{"r", pi.Variable("rresult"), pi.Send{"p", pi.Sub{pi.Variable("lresult"), pi.Variable("rresult")}}}}}}}
	case lambda.Mult:
		return pi.Privatize{"l", pi.Privatize{"r", pi.Parallel{Translate(v.L, "l"), Translate(v.R, "r"), pi.ReceiveThen{"l", pi.Variable("lresult"), pi.ReceiveThen{"r", pi.Variable("rresult"), pi.Send{"p", pi.Mul{pi.Variable("lresult"), pi.Variable("rresult")}}}}}}}
	case lambda.Div:
		return pi.Privatize{"l", pi.Privatize{"r", pi.Parallel{Translate(v.L, "l"), Translate(v.R, "r"), pi.ReceiveThen{"l", pi.Variable("lresult"), pi.ReceiveThen{"r", pi.Variable("rresult"), pi.Send{"p", pi.Div{pi.Variable("lresult"), pi.Variable("rresult")}}}}}}}
	default:
		panic("not supposed to happen")
	}
}

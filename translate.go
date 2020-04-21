package main

import (
	"pingo/lambda"
)

func translate(lexpr lambda.Lambda, channel string) expr {
	switch v := lexpr.(type) {
	case lambda.Lconst:
		return send{channel, constant(v)}
	case lambda.Lvar:
		return send{channel, variable(v)}
	case lambda.Lfun:
		return privatize{"y", parallel{send{channel, variable("y")}, receiveThen{"y", pair{variable(v.Arg), variable("q")}, translate(v.Exp, "q")}}}
	case lambda.Lapp:
		return privatize{"q", parallel{translate(v.Exp, "q"), receiveThen{"q", variable("v"), privatize{"r", parallel{translate(v.Fun, "r"), receiveThen{"r", variable("f"), send{"f", pair{variable("v"), variable("p")}}}}}}}}
	default:
		panic("not supposed to happen")
	}
}
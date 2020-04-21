//go:generate goyacc -o src/pi/pilang.go -v src/pi/pilang.output -p "expr" src/pi/pilang.y
//go:generate goyacc -o src/lambda/lambda.go -v src/lambda/lambda.output -p "lambda" src/lambda/lambda.y

package main

import (
	"flag"
	"fmt"
	"pingo/src/lambda"
	"pingo/src/pi"
	"pingo/src/translate"
	"sync"
)

func main() {
	// The options of the executable
	showSrc := flag.Bool("showsrc", false, "Output the parsed code")
	outCode := flag.Bool("outcode", false, "Output the code before executing it (after translation if any")
	translateInput := flag.Bool("translate", false, "Parse input as lambda, translate in pi and execute")
	flag.Parse()

	var ret pi.Expr
	if *translateInput {
		if *showSrc {
			lambda.Test()
		}
		t := translate.Translate(lambda.GetParsedInput(), "p")
		ret = pi.Parallel{t, pi.ReceiveThen{"p", pi.Variable("x"), pi.Print{pi.Variable("x"), pi.Skip(0)}}}
	} else {
		if *showSrc {
			pi.Test()
		}
		ret = pi.GetParsedInput()
	}

	if ret == nil {
		return
	}

	// We print the code if asked
	if *outCode {
		fmt.Println(ret)
	}

	// We want all goroutines to have finished before killing the process
	var wg sync.WaitGroup

	wg.Add(1)
	// We launch the evaluation in a goroutine
	go pi.Launch(ret, &wg)

	wg.Wait() // We wait for all goroutines to terminate
}

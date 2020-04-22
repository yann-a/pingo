//go:generate goyacc -o src/pi/pilang.go -v src/pi/pilang.output -p "expr" src/pi/pilang.y
//go:generate goyacc -o src/lambda/lambda.go -v src/lambda/lambda.output -p "lambda" src/lambda/lambda.y

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

	// Parsing
	in := getBuffer()
	ret := Parse(in, *translateInput, *showSrc)

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

// returns parsed input from buffer
func Parse(in *bufio.Reader, translateInput, showSrc bool) (ret pi.Expr) {
	if translateInput {
		lambdaCode := lambda.Parse(in)
		if showSrc {
			fmt.Println(lambdaCode)
		}
		t := translate.Translate(lambdaCode, "p")
		ret = pi.Parallel{t, pi.ReceiveThen{"p", pi.Variable("x"), pi.Print{pi.Variable("x"), pi.Skip(0)}}}
	} else {
		ret = pi.Parse(in)
		if showSrc {
			fmt.Println(ret)
		}
	}
	return
}

// returns the right buffer according to options
func getBuffer() *bufio.Reader {
	nonFlagArgs := flag.Args()
	buffer := os.Stdin

	// If a file is provided we try reading from it
	if len(nonFlagArgs) > 0 {
		file, err := os.Open(nonFlagArgs[0])
		if err != nil {
			fmt.Printf("Couldn't read from %s (%s). Reading from stdin\n", nonFlagArgs[0], err)
		} else {
			defer file.Close()
			buffer = file
		}
	}

	return bufio.NewReader(buffer)
}
